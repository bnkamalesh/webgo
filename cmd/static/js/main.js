const webgo = async () => {
  const clientID = Math.random()
    .toString(36)
    .replace(/[^a-z]+/g, "")
    .substring(0, 16);  
  const sseDOM = document.getElementById("sse");

  const startSSE = (clientID) => {
    const source = new EventSource(`/sse/${clientID}`);
    source.onmessage = function (event) {
      const time = new Date(event.data);
      sseDOM.innerText = `[${clientID}]: ${time.toLocaleString()}`;
    };
    source.onerror = (err) => {
      source.close()
      console.log(err);
      sseDOM.innerText = `SSE error, restarting`;
      startSSE(clientID)
    };
  };
  startSSE(clientID)
};
webgo();
