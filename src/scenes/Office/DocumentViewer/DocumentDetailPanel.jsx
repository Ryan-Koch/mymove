import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, omit } from 'lodash';
import { reduxForm, getFormValues, isValid, FormSection } from 'redux-form';

import editablePanel from '../editablePanel';
import { renderStatusIcon } from 'shared/utils';
import { formatDate, formatCents } from 'shared/formatters';
import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import ExpenseDocumentForm from './ExpenseDocumentForm';
import {
  selectMoveDocument,
  updateMoveDocument,
} from 'shared/Entities/modules/moveDocuments';
import { isMovingExpenseDocument } from 'shared/Entities/modules/movingExpenseDocuments';

import '../office.css';

const DocumentDetailDisplay = props => {
  const moveDoc = props.moveDocument;
  const isExpenseDocument = isMovingExpenseDocument(moveDoc);
  const moveDocFieldProps = {
    values: moveDoc,
    schema: props.moveDocSchema,
  };
  const reimbursementFieldProps = {
    values: get(moveDoc, 'reimbursement', {}),
    schema: props.reimbursementSchema,
  };
  return (
    <React.Fragment>
      <div>
        <span className="panel-subhead">
          {renderStatusIcon(moveDoc.status)}
          {moveDoc.title}
        </span>
        <p className="uploaded-at">
          Uploaded {formatDate(get(moveDoc, 'document.uploads.0.created_at'))}
        </p>
        {moveDoc.title ? (
          <PanelSwaggerField fieldName="title" {...moveDocFieldProps} />
        ) : (
          <PanelField title="Document Title" className="missing">
            Missing
          </PanelField>
        )}
        {moveDoc.move_document_type ? (
          <PanelSwaggerField
            fieldName="move_document_type"
            {...moveDocFieldProps}
          />
        ) : (
          <PanelField title="Document Type" className="missing">
            Missing
          </PanelField>
        )}
        {isExpenseDocument &&
          moveDoc.moving_expense_type && (
            <PanelSwaggerField
              fieldName="moving_expense_type"
              {...moveDocFieldProps}
            />
          )}
        {isExpenseDocument &&
          get(moveDoc, 'reimbursement.requested_amount') && (
            <PanelSwaggerField
              fieldName="requested_amount"
              {...reimbursementFieldProps}
            />
          )}
        {isExpenseDocument &&
          get(moveDoc, 'reimbursement.method_of_receipt') && (
            <PanelSwaggerField
              title="Payment Method"
              fieldName="method_of_receipt"
              {...reimbursementFieldProps}
            />
          )}
        {moveDoc.status ? (
          <PanelSwaggerField fieldName="status" {...moveDocFieldProps} />
        ) : (
          <PanelField title="Document Status" className="missing">
            Missing
          </PanelField>
        )}
        {moveDoc.notes ? (
          <PanelSwaggerField fieldName="notes" {...moveDocFieldProps} />
        ) : (
          <PanelField title="Notes" />
        )}
      </div>
    </React.Fragment>
  );
};

const DocumentDetailEdit = props => {
  const { formValues, moveDocSchema, reimbursementSchema } = props;
  const isExpenseDocument =
    get(formValues, 'moveDocument.move_document_type', '') === 'EXPENSE';
  return (
    <React.Fragment>
      <div>
        <FormSection name="moveDocument">
          <SwaggerField fieldName="title" swagger={moveDocSchema} required />
          <SwaggerField
            fieldName="move_document_type"
            swagger={moveDocSchema}
            required
          />
          {isExpenseDocument && (
            <ExpenseDocumentForm
              moveDocSchema={moveDocSchema}
              reimbursementSchema={reimbursementSchema}
            />
          )}
          <SwaggerField fieldName="status" swagger={moveDocSchema} required />
          <SwaggerField fieldName="notes" swagger={moveDocSchema} />
        </FormSection>
      </div>
    </React.Fragment>
  );
};

const formName = 'move_document_viewer';

let DocumentDetailPanel = editablePanel(
  DocumentDetailDisplay,
  DocumentDetailEdit,
);
DocumentDetailPanel = reduxForm({ form: formName })(DocumentDetailPanel);

function mapStateToProps(state, props) {
  const moveDocumentId = props.moveDocumentId;
  let moveDocument = selectMoveDocument(state, moveDocumentId);
  // Don't pass 0-value reimbursement values to update endpoint
  if (
    get(moveDocument, 'reimbursement.id', '') ===
    '00000000-0000-0000-0000-000000000000'
  ) {
    moveDocument = omit(moveDocument, 'reimbursement');
  }
  // Convert cents to collars - make a deep clone copy to not modify moveDocument itself
  let initialMoveDocument = JSON.parse(JSON.stringify(moveDocument));
  let requested_amount = get(
    initialMoveDocument,
    'reimbursement.requested_amount',
  );
  if (requested_amount) {
    initialMoveDocument.reimbursement.requested_amount = formatCents(
      requested_amount,
    );
  }

  return {
    // reduxForm
    initialValues: {
      moveDocument: initialMoveDocument,
    },
    formValues: getFormValues(formName)(state),
    moveDocSchema: get(
      state,
      'swagger.spec.definitions.MoveDocumentPayload',
      {},
    ),
    reimbursementSchema: get(
      state,
      'swagger.spec.definitions.Reimbursement',
      {},
    ),
    hasError: false,
    errorMessage: state.office.error,
    isUpdating: false,
    moveDocument: moveDocument,

    // editablePanel
    formIsValid: isValid(formName)(state),
    getUpdateArgs: function() {
      // Make a copy of values to not modify moveDocument
      let values = JSON.parse(JSON.stringify(getFormValues(formName)(state)));
      values.moveDocument.personally_procured_move_id = get(
        state.office,
        'officePPMs.0.id',
      );
      if (
        get(values.moveDocument, 'move_document_type', '') !== 'EXPENSE' &&
        get(values.moveDocument, 'reimbursement', false)
      ) {
        values.moveDocument = omit(values.moveDocument, 'reimbursement');
      }
      if (get(values.moveDocument, 'move_document_type', '') === 'EXPENSE') {
        values.moveDocument.reimbursement.requested_amount = parseFloat(
          values.moveDocument.reimbursement.requested_amount,
        );
      }
      let requested_amount = get(
        values.moveDocument,
        'reimbursement.requested_amount',
        '',
      );
      if (requested_amount) {
        values.moveDocument.reimbursement.requested_amount = Math.round(
          parseFloat(requested_amount) * 100,
        );
      }
      return [
        get(state, 'office.officeMove.id'),
        get(moveDocument, 'id'),
        values.moveDocument,
      ];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: updateMoveDocument,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(
  DocumentDetailPanel,
);
