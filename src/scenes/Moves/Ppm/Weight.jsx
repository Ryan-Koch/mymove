import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';
import PropTypes from 'prop-types';
import Slider from 'react-rangeslider'; //todo: pull from node_modules, override

import WizardPage from 'shared/WizardPage';
import Alert from 'shared/Alert';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import {
  setPendingPpmWeight,
  loadPpm,
  getPpmWeightEstimate,
  createOrUpdatePpm,
} from './ducks';

import 'react-rangeslider/lib/index.css';
import './Weight.css';

function getWeightInfo(ppm) {
  const size = ppm ? ppm.size : 'L';
  switch (size) {
    case 'S':
      return {
        min: 100,
        max: 800,
      };
    case 'M':
      return {
        min: 400,
        max: 1200,
      };
    default:
      return {
        min: 1000,
        max: 5000, //TODO: this should be max entitlement
      };
  }
}
export class PpmWeight extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Weight Selection';
    if (!this.props.currentPpm) {
      this.props.loadPpm(this.props.match.params.moveId);
    } else {
      this.updateIncentive();
    }
  }
  componentDidUpdate(prevProps, prevState) {
    if (!prevProps.hasLoadSuccess && this.props.hasLoadSuccess) {
      this.updateIncentive();
    }
  }
  // this method is used to set the incentive on page load
  // it runs even if the incentive has been set before since data changes on previous pages could
  // affect it
  updateIncentive() {
    const { pendingPpmWeight, currentWeight, currentPpm } = this.props;
    const weight_estimate = get(this.props, 'currentPpm.weight_estimate');
    if (![pendingPpmWeight, weight_estimate].includes(currentWeight)) {
      this.onWeightSelecting(currentWeight);
      this.props.getPpmWeightEstimate(
        currentPpm.planned_move_date,
        currentPpm.pickup_zip,
        currentPpm.destination_zip,
        currentWeight,
      );
    }
  }
  handleSubmit = () => {
    const { pendingPpmWeight, incentive, createOrUpdatePpm } = this.props;
    const moveId = this.props.match.params.moveId;
    createOrUpdatePpm(moveId, {
      weight_estimate: pendingPpmWeight,
      estimated_incentive: incentive,
    });
  };
  onWeightSelecting = value => {
    this.props.setPendingPpmWeight(value);
  };
  onWeightSelected = value => {
    const { currentPpm } = this.props;
    this.props.getPpmWeightEstimate(
      currentPpm.planned_move_date,
      currentPpm.pickup_zip,
      currentPpm.destination_zip,
      this.props.pendingPpmWeight,
    );
  };
  render() {
    const {
      pendingPpmWeight,
      currentPpm,
      incentive,
      pages,
      pageKey,
      hasSubmitSuccess,
      currentWeight,
      hasLoadSuccess,
      error,
    } = this.props;
    const currentInfo = getWeightInfo(currentPpm);

    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={true}
        pageIsDirty={Boolean(pendingPpmWeight)}
        hasSucceeded={hasSubmitSuccess}
      >
        {error && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {error.message}
            </Alert>
          </div>
        )}
        <h2>Customize Weight</h2>
        {!hasLoadSuccess && <LoadingPlaceholder />}
        {hasLoadSuccess && (
          <React.Fragment>
            <p>
              Use this slider to customize how much weight you think you’ll
              carry.
            </p>
            <div className="slider-container">
              <Slider
                min={currentInfo.min}
                max={currentInfo.max}
                value={currentWeight}
                onChange={this.onWeightSelecting}
                onChangeComplete={this.onWeightSelected}
                labels={{
                  [currentInfo.min]: currentInfo.min,
                  [currentInfo.max]: currentInfo.max,
                  //[currentWeight]: currentWeight,
                }}
              />
            </div>
            <table className="numeric-info">
              <tr>
                <th>Your PPM Weight Estimate:</th>
                <td className="current-weight"> {currentWeight}</td>
              </tr>
              <tr>
                <th>Your PPM Incentive:</th>
                <td className="incentive">{incentive}</td>
              </tr>
            </table>

            <div className="info">
              <h3> How is my PPM Incentive calculated?</h3>
              <p>
                The government gives you 95% of what they would pay a mover when
                you move your own belongings, based on weight and distance. You
                pay taxes on this income. You can reduce the amount taxable
                incentive by saving receipts for approved expenses.
              </p>

              <p>
                This estimator just presents a range of possible incentives
                based on your anticipated shipment weight, anticipated moving
                date, and the specific route that you will be traveling. During
                your move, you will need to weigh the stuff you’re carrying, and
                submit weight tickets. We’ll let you know later how to weigh the
                stuff you carry.
              </p>
            </div>
          </React.Fragment>
        )}
      </WizardPage>
    );
  }
}

PpmWeight.propTypes = {
  pendingPpmWeight: PropTypes.number,
  currentWeight: PropTypes.number,
  currentPpm: PropTypes.shape({
    id: PropTypes.string,
    size: PropTypes.string,
    weight: PropTypes.number,
    incentive: PropTypes.string,
  }),
  hasSubmitSuccess: PropTypes.bool.isRequired,
  hasLoadSuccess: PropTypes.bool.isRequired,
  setPendingPpmWeight: PropTypes.func.isRequired,
};

function getMiddleWeight(ppm) {
  const currentInfo = getWeightInfo(ppm);
  return currentInfo.min + (currentInfo.max - currentInfo.min) / 2;
}
function mapStateToProps(state) {
  const currentWeight =
    state.ppm.pendingPpmWeight ||
    get(
      state,
      'ppm.currentPpm.weight_estimate',
      getMiddleWeight(state.ppm.currentPpm),
    );
  const props = {
    ...state.ppm,
    currentWeight,
  };

  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { setPendingPpmWeight, loadPpm, getPpmWeightEstimate, createOrUpdatePpm },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(PpmWeight);
