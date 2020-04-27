import React from "react";

const FileUploader = () => {
  const handleFile = (event) => {
    console.log(event.target.files[0]);
  };

  return (
    <>
      <div className="fileContainer">
        <p>Selecione um arquivo</p>
        <input type="file" onChange={handleFile} />
      </div>

      <button className="uploadBtn">Enviar</button>
    </>
  );
};

export default FileUploader;
