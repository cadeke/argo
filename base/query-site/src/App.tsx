import { useState } from "react";
import IpToDomain from "./IpToDomain";
import DomainToIp from "./DomainToIp";
import "./App.css";

function App() {
  const [view, setView] = useState("ip2domain");

  return (
    <div className="container">
      <div className="container_menu">
        <h1>ARGO Query Site</h1>

        <div className="header" style={{ marginBottom: "1rem" }}>
          <button
            className="IpToDomainBtn"
            onClick={() => setView("ip2domain")}
          >
            <i className="material-icons">public</i> IP to Domain
          </button>
          <button
            className="DomainToIpBtn"
            onClick={() => setView("domain2ip")}
          >
            <i className="material-icons">dns</i> Domain to IP
          </button>
        </div>
        <div className="content">
          {view === "ip2domain" && <IpToDomain />}
          {view === "domain2ip" && <DomainToIp />}
        </div>
      </div>
    </div>
  );
}

export default App;
