import { useState, useEffect } from "react";

function GetData() {
  const [data, setData] = useState([]);

  const host = import.meta.env.VITE_API_HOST;
  const port = import.meta.env.VITE_API_PORT;

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(`http://${host}:${port}/api/list`);
        const result = await response.json();

        console.log("Data fetched:", result);

        setData(result);
      } catch (error) {
        console.error("Error fetching data:", error);
      }
    };

    fetchData();
  }, []);

  return (
    <div>
      <h2>GET All Records</h2>
      <div className="form-container">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Domain</th>
              <th>IP</th>
              <th>Server</th>
            </tr>
          </thead>
          <tbody>
            {data.map((item: any, index) => (
              <tr key={index}>
                <td>{item.id}</td>
                <td>{item.domain}</td>
                <td>{item.ip}</td>
                <td>{item.server}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default GetData;
