import { useState } from "react";

function DeleteData() {
  const [id, setId] = useState("");

  const host = import.meta.env.VITE_API_HOST;
  const port = import.meta.env.VITE_API_PORT;

  const handleDelete = async () => {
    try {
      const response = await fetch(
        `http://${host}:${port}/api/delete?id=${id}`,
        {
          method: "DELETE",
        }
      );
      const result = await response.json();
      console.log("Data deleted:", result);
      alert(`Data successfully deleted with record ID ${id}`);
    } catch (error) {
      console.error("Error deleting data:", error);
    }
  };

  return (
    <div>
      <h2>DELETE Record</h2>
      <div className="form-container">
        <input
          type="text"
          placeholder="Enter ID"
          value={id}
          onChange={(e) => setId(e.target.value)}
        />
        <button className="deleteBtn" onClick={handleDelete}>
          Delete
        </button>
      </div>
    </div>
  );
}

export default DeleteData;
