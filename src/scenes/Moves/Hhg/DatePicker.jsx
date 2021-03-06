import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { DayPicker } from 'react-day-picker';
import { withRouter } from 'react-router-dom';
import { get } from 'lodash';
import 'react-day-picker/lib/style.css';
import moment from 'moment';

import { formatSwaggerDate, parseSwaggerDate } from 'shared/formatters';
import { bindActionCreators } from 'redux';
import { getMoveDatesSummary, selectMoveDatesSummary } from 'shared/Entities/modules/moves';
import DatesSummary from 'scenes/Moves/Hhg/DatesSummary.jsx';

import './DatePicker.css';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { loadEntitlementsFromState } from 'shared/entitlements';

const getRequestLabel = 'DatePicker.getMoveDatesSummary';

function createModifiers(moveDates) {
  if (!moveDates) {
    return null;
  }

  return {
    pack: convertDateStringArray(moveDates.pack),
    pickup: convertDateStringArray(moveDates.pickup),
    transit: convertDateStringArray(moveDates.transit),
    delivery: convertDateStringArray(moveDates.delivery),
    report: convertDateStringArray(moveDates.report),
  };
}

function convertDateStringArray(dateStrings) {
  return dateStrings && dateStrings.map(dateString => parseSwaggerDate(dateString));
}

export class HHGDatePicker extends Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedDay: null,
    };
  }

  handleDayClick = (day, { disabled }) => {
    if (disabled) {
      return;
    }
    const moveDate = formatSwaggerDate(day);
    this.props.input.onChange(moveDate);
    this.props.getMoveDatesSummary(getRequestLabel, this.props.match.params.moveId, moveDate);
    this.setState({
      selectedDay: moveDate,
    });
  };

  isDayDisabled = day => {
    const availableMoveDates = this.props.availableMoveDates;
    if (!availableMoveDates) {
      return true;
    }

    const momentDay = moment(day);
    if (momentDay.isBefore(availableMoveDates.minDate, 'day') || momentDay.isAfter(availableMoveDates.maxDate, 'day')) {
      return true;
    }

    return !availableMoveDates.available.find(element => momentDay.isSame(element, 'day'));
  };

  componentDidUpdate(prevProps) {
    const moveDateIsSavedDate =
      get(this.props, 'currentShipment.requested_pickup_date') &&
      this.props.currentShipment.requested_pickup_date === this.props.input.value;
    if (this.props.currentShipment !== prevProps.currentShipment && this.props.currentShipment.requested_pickup_date) {
      if (!moveDateIsSavedDate) {
        this.props.getMoveDatesSummary(
          getRequestLabel,
          this.props.match.params.moveId,
          this.props.currentShipment.requested_pickup_date,
        );
      }
      this.setState({
        selectedDay: this.props.input.value || this.props.currentShipment.requested_pickup_date,
      });
    }
  }

  componentDidMount() {
    this.setState({
      selectedDay: this.props.input.value || get(this.props, 'currentShipment.requested_pickup_date'),
    });
  }

  render() {
    const availableMoveDates = this.props.availableMoveDates;
    const parsedSelectedDay = parseSwaggerDate(this.state.selectedDay);
    const entitlementSum = get(this.props, 'entitlement.sum');
    return (
      <div className="form-section">
        {availableMoveDates ? (
          <div className="usa-grid">
            <h3 className="instruction-heading">Pick a moving date.</h3>
            <div className="usa-width-one-third">
              <DayPicker
                onDayClick={this.handleDayClick}
                month={parsedSelectedDay || (availableMoveDates && availableMoveDates.minDate)}
                selectedDays={parsedSelectedDay}
                disabledDays={this.isDayDisabled}
                modifiers={this.props.modifiers}
                showOutsideDays
              />
            </div>
            <div className="usa-width-two-thirds">
              {this.state.selectedDay && <DatesSummary moveDates={this.props.moveDates} />}
            </div>
          </div>
        ) : (
          <LoadingPlaceholder />
        )}
        <div className="usa-grid pack-days-notice">
          <div className="usa-width-one-whole">
            <p>Can't find a date that works? Talk with a move counselor in your local Transportation office (PPPO).</p>
            {entitlementSum ? (
              <div className="weight-info-box">
                * It takes 1 day for every 5,000 lbs of stuff movers need to pack. You have an allowance of{' '}
                {entitlementSum.toLocaleString()} lbs, so we estimate it will take {Math.ceil(entitlementSum / 5000)}{' '}
                days to pack.
              </div>
            ) : (
              <LoadingPlaceholder />
            )}
          </div>
        </div>
      </div>
    );
  }
}

HHGDatePicker.propTypes = {
  input: PropTypes.object.isRequired,
  currentShipment: PropTypes.object,
  availableMoveDates: PropTypes.object,
};

function mapStateToProps(state, ownProps) {
  const moveDate = ownProps.input.value || get(ownProps, 'currentShipment.requested_pickup_date');
  const moveDateIsSavedDate =
    get(ownProps, 'currentShipment.requested_pickup_date') &&
    ownProps.currentShipment.requested_pickup_date === ownProps.input.value;
  const moveDates = moveDateIsSavedDate
    ? ownProps.currentShipment.move_dates_summary
    : selectMoveDatesSummary(state, ownProps.match.params.moveId, moveDate);
  return {
    moveDates: moveDates,
    modifiers: createModifiers(moveDates),
    entitlement: loadEntitlementsFromState(state),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ getMoveDatesSummary }, dispatch);
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(HHGDatePicker));
