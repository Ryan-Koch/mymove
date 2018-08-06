import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import { NavLink } from 'react-router-dom';

import { withContext } from 'shared/AppContext';

class ShipmentInfo extends Component {
  render() {
    return (
      <div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-two-thirds">
            <h1>Shipment Info: LastName, FirstName</h1>
          </div>
          <div className="usa-width-one-third nav-controls">
            <NavLink to="/queues/new" activeClassName="usa-current">
              <span>New Shipments Queue</span>
            </NavLink>
          </div>
        </div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-one-whole">
            <ul className="move-info-header-meta Todo-phase2">
              <li>GBL# KKFA9999999</li>
              <li>Locator# ABC89</li>
              <li>KKFA to HAFC</li>
              <li>Move date 07-Jun-2018</li>
              <li>
                Status: <b>At Destination</b>
              </li>
            </ul>
          </div>
        </div>
        <div className="usa-grid grid-wide tabs">
          <div className="usa-width-two-thirds">
            <p>
              <button className="usa-button-primary">Accept</button>
              <button className="usa-button-secondary">Reject</button>
            </p>
          </div>
          <div className="usa-width-one-third" />
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => ({});

const mapDispatchToProps = dispatch => bindActionCreators({}, dispatch);

export default withContext(
  connect(mapStateToProps, mapDispatchToProps)(ShipmentInfo),
);
