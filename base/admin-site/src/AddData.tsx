import { useState } from "react";

function AddData() {
  const [formData, setFormData] = useState({});

  const host = import.meta.env.VITE_API_HOST;
  const port = import.meta.env.VITE_API_PORT;

  const handleChange = (e: { target: { name: any; value: any } }) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: { preventDefault: () => void }) => {
    e.preventDefault();
    try {
      const response = await fetch(`http://${host}:${port}/api/add`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });
      const result = await response.json();

      if (response.status == 201) {
        console.log("Data added:", result);
        alert(`Data successfully added with record ID ${result.id}`);
      } else {
        console.error("Error adding data:", result);
        alert("Error adding data" + result.message);
      }
    } catch (error) {
      console.error("Error adding data:", error);
    }
  };

  return (
    <div className="data-container">
      <h2>ADD Record</h2>
      <div className="form-container">
        <form onSubmit={handleSubmit}>
          <label className="form-label">IP:</label>
          <input
            type="text"
            name="ip"
            placeholder="Input IP address here..."
            onChange={handleChange}
          />
          <br />
          <label className="form-label">Domain:</label>
          <input
            type="text"
            name="domain"
            placeholder="Input domain name here..."
            onChange={handleChange}
          />
          <br />
          <label className="form-label">Server:</label>
          <input
            type="text"
            name="server"
            placeholder="Input server IP here..."
            onChange={handleChange}
          />
          <br />
          <button className="addBtn" type="submit">
            Add
          </button>
        </form>
      </div>
    </div>
  );
}

export default AddData;
