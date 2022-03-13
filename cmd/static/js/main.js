const webgo = async () => {
  const clientID = Math.random()
    .toString(36)
    .replace(/[^a-z]+/g, "")
    .substring(0, 16);
  const sseDOM = document.getElementById("sse");
  const sseClientsDOM = document.getElementById("sse-clients");
  const sseClientIDDOM = document.getElementById("sse-client-id");

  const startSSE = (clientID) => {
    const source = new EventSource(`/sse/${clientID}`);
    source.onmessage = function (event) {
      const parts = event.data?.split("(");
      const date = new Date(parts[0]);
      const activeClients = parts[1].replace(")", "");
      sseDOM.innerText = date.toLocaleString();
      sseClientsDOM.innerText = activeClients;
      sseClientIDDOM.innerText = clientID;
    };
    source.onerror = (err) => {
      source.close();
      console.log(err);
      sseDOM.innerText = `SSE error, restarting`;
      startSSE(clientID);
    };
  };
  startSSE(clientID);
};
webgo();
