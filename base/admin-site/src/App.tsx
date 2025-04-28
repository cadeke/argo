import { useState } from "react";
import GetData from "./GetData";
import AddData from "./AddData";
import UpdateData from "./UpdateData";
import DeleteData from "./DeleteData";
import "./App.css";

function App() {
  const [view, setView] = useState("get");

  return (
    <div className="container">
      <div className="container_menu">
        <h1>ARGO Admin Site</h1>

        <div className="header">
          <button className="getBtn" onClick={() => setView("get")}>
            <i className="material-icons">search</i> GET
          </button>
          <button className="addBtn" onClick={() => setView("add")}>
            <i className="material-icons">add</i> ADD
          </button>
          <button className="updateBtn" onClick={() => setView("update")}>
            <i className="material-icons">edit</i> UPDATE
          </button>
          <button className="deleteBtn" onClick={() => setView("delete")}>
            <i className="material-icons">delete</i> DELETE
          </button>
        </div>

        <div className="content">
          {view === "get" && <GetData />}
          {view === "add" && <AddData />}
          {view === "update" && <UpdateData />}
          {view === "delete" && <DeleteData />}
        </div>
      </div>
    </div>
  );
}

export default App;
