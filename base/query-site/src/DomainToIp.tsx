import { useState } from "react";
function DomainToIp() {
  const [domain, setDomain] = useState("");
  const [result, setResult] = useState("");

  const handleTransform = async () => {
    const host = import.meta.env.VITE_API_HOST;
    const port = import.meta.env.VITE_API_PORT;    

    try {
      const response = await fetch(
        `http://${host}:${port}/api/domain2ip?query=${domain}`,
      );
      const data = await response.json();
      console.log("Data:", data);

      if (response.status !== 200) {
        setResult("Error resolving IP");
        return;
      }

      setResult(`IP for ${domain}: ${data.result}`);
    } catch (error) {
      console.error("Error fetching data:", error);
      setResult("Error resolving domain");
    }
  };

  return (
    <div>
      <h2>Domain to IP</h2>
      <div className="form-container">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleTransform();
          }}
        >
          <input
            type="text"
            placeholder="Enter domain"
            value={domain}
            onChange={(e) => setDomain(e.target.value)}
          />
          <button className="DomainToIpBtn" type="submit">
            Go
          </button>
        </form>
        {result && <p>{result}</p>}
      </div>
    </div>
  );
}

export default DomainToIp;
