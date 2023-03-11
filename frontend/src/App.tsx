import axios from "axios";
import { useEffect, useState } from "react";
import reactLogo from "./assets/react.svg";
import "./App.css";
// axios.defaults.xsrfCookieName = "_csrf";
// axios.defaults.xsrfHeaderName = "X-CSRF-TOKEN";
axios.defaults.withCredentials = true;

function App() {
  const [count, setCount] = useState(0);

  useEffect(() => {
    async function getToken() {
      const res = await axios.get("http://localhost:8080/api/csrf");
      console.log(res.headers);
      axios.defaults.headers.common["X-CSRF-Token"] =
        res.headers["x-csrf-token"];
    }
    getToken();
  }, []);

  const handleSubmit = async () => {
    const res = await axios.post("http://localhost:8080/api/signup");
    console.log(res);
  };

  return (
    <div className="App">
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src="/vite.svg" className="logo" alt="Vite logo" />
        </a>
        <a href="https://reactjs.org" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div>
        <input type="text" name="email" />
        <button onClick={handleSubmit}>Submit</button>
      </div>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </div>
  );
}

export default App;
