import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';
import editablePanel from './editablePanel';

import { updateAccounting } from './ducks';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelSwaggerField } from 'shared/EditablePanel';

const AccountingDisplay = props => {
  const fieldProps = {
    schema: props.ordersSchema,
    values: props.orders,
  };

  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelSwaggerField fieldName="dept_indicator" {...fieldProps} />
      </div>
      <div className="editable-panel-column">
        <PanelSwaggerField fieldName="tac" {...fieldProps} />
      </div>
    </React.Fragment>
  );
};

const AccountingEdit = props => {
  const { schema } = props;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <SwaggerField fieldName="dept_indicator" swagger={schema} required />
      </div>
      <div className="editable-panel-column">
        <SwaggerField fieldName="tac" swagger={schema} required />
      </div>
    </React.Fragment>
  );
};

const formName = 'office_move_info_accounting';

let AccountingPanel = editablePanel(AccountingDisplay, AccountingEdit);
AccountingPanel = reduxForm({ form: formName })(AccountingPanel);

function mapStateToProps(state) {
  return {
    // reduxForm
    formData: state.form[formName],
    initialValues: state.office.accounting,

    // Wrapper
    ordersSchema: get(state, 'swagger.spec.definitions.PatchAccounting', {}),
    hasError:
      state.office.accountingHasLoadError ||
      state.office.accountingHasUpdateError,
    errorMessage: state.office.error,
    orders: state.office.accounting || {},
    isUpdating: state.office.accountingIsUpdating,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: updateAccounting,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(AccountingPanel);
