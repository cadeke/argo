import { useState } from "react";

function IpToDomain() {
  const [ip, setIp] = useState("");
  const [result, setResult] = useState("");

  const handleTransform = async () => {
    const host = import.meta.env.VITE_API_HOST;
    const port = import.meta.env.VITE_API_PORT;

    try {
      const response = await fetch(
        `http://${host}:${port}/api/ip2domain?query=${ip}`
      );
      const data = await response.json();
      console.log("Data:", data);

      if (response.status !== 200) {
        setResult("Error resolving domain");
        return;
      }

      setResult(`Domain for ${ip}: ${data.result}`);
    } catch (error) {
      console.error("Error fetching data:", error);
      setResult("Error resolving domain");
    }
  };

  return (
    <div>
      <h2>IP to Domain</h2>
      <div className="form-container">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleTransform();
          }}
        >
          <input
            type="text"
            placeholder="Enter IP"
            value={ip}
            onChange={(e) => setIp(e.target.value)}
          />
          <button className="IpToDomainBtn" type="submit">
            Go
          </button>
        </form>
        {result && <p>{result}</p>}
      </div>
    </div>
  );
}

export default IpToDomain;
