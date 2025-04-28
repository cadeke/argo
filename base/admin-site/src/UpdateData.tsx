import { useState } from "react";

function UpdateData() {
  const [formData, setFormData] = useState({
    id: "",
    domain: "",
    ip: "",
    server: "",
  });

  const host = import.meta.env.VITE_API_HOST;
  const port = import.meta.env.VITE_API_PORT;

  const handleChange = (e: { target: { name: string; value: string } }) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: { preventDefault: () => void }) => {
    e.preventDefault();
    try {
      const response = await fetch(`http://${host}:${port}/api/update`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });
      const result = await response.json();
      console.log("Data updated:", result);
      alert(`Data successfully updated with record ID ${result.id}`);
    } catch (error) {
      console.error("Error updating data:", error);
    }
  };

  return (
    <div>
      <h2>UPDATE Record</h2>
      <div className="form-container">
        <form onSubmit={handleSubmit}>
          <label className="form-label">ID:</label>
          <input
            type="text"
            name="id"
            placeholder="Input ID here..."
            value={formData.id}
            onChange={handleChange}
          />
          <br />
          <label className="form-label">Domain:</label>
          <input
            type="text"
            name="domain"
            placeholder="Input domain here..."
            value={formData.domain}
            onChange={handleChange}
          />
          <br />
          <label className="form-label">IP:</label>
          <input
            type="text"
            name="ip"
            placeholder="Input IP address here..."
            value={formData.ip}
            onChange={handleChange}
          />
          <br />
          <label className="form-label">Server:</label>
          <input
            type="text"
            name="server"
            placeholder="Input server here..."
            value={formData.server}
            onChange={handleChange}
          />
          <br />
          <button className="updateBtn" type="submit">
            Update
          </button>
        </form>
      </div>
    </div>
  );
}

export default UpdateData;
