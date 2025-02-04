import React from 'react';
import '../index';
import { BoldBI } from '@boldbi/boldbi-embedded-sdk';

//Url of the authorizationserver action in the Go application(http://localhost:8086/authorizationserver). Learn more about authorize server [here](https://help.syncfusion.com/bold-bi/embedded-bi/javascript/authorize-server)
const authorizationUrl = "http://localhost:8086/authorizationServer";

class Dashboard extends React.Component {
  constructor(props) {
    super(props);
    this.state = { toke: undefined, items: [] };
    this.BoldBiObj = new BoldBI();
  };

  renderDashboard(embedConfig) {
    this.dashboard = BoldBI.create({
      serverUrl: embedConfig.ServerUrl + "/" + embedConfig.SiteIdentifier,
      dashboardId: embedConfig.DashboardId,
      embedContainerId: "dashboard",
      width: "100%",
      height: window.innerHeight + 'px',
      authorizationServer: {
        url: authorizationUrl
      }
    });
    this.dashboard.loadDashboard();
  }

  render() {
    return (
      <div id="DashboardListing">
        <div id="viewer-section">
          <div id="dashboard"></div>
        </div>
      </div>
    );
  }

  async componentDidMount() {
    try {
      const response = await fetch('http://localhost:8086/getServerDetails');
      const embedConfig = await response.json();
      this.renderDashboard(embedConfig);  
    } catch (error) {
      console.log("Error: embedConfig.json file not found.");
      this.setState({ toke: "error", items: "error" });
    }
  }
}
export default Dashboard;