import React from "react";
import logo from "./logo.svg";
import "./App.css";
import FileUploader from "./components/FileUploader";

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <div>
          <FileUploader />
        </div>
      </header>
    </div>
  );
}

export default App;
