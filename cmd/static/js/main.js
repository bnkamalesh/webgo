const webgo = async () => {
  const clientID = Math.random()
    .toString(36)
    .replace(/[^a-z]+/g, "")
    .substring(0, 16);
  const sseDOM = document.getElementById("sse");
  const sseClientsDOM = document.getElementById("sse-clients");
  const sseClientIDDOM = document.getElementById("sse-client-id");

  const formatBackoff = (backoff, precision = 2) => {
    let boff = `${backoff}ms`;
    if (backoff >= 1000) {
      boff = `${parseFloat(backoff / 1000).toFixed(precision)}s`;
    }
    return boff;
  };

  const config = {
    url: `/sse/${clientID}`,
    onMessage: (data) => {
      const parts = data?.split?.("(");
      if (!parts || !parts.length) {
        return;
      }

      const date = new Date(parts[0]);
      const activeClients = parts[1].replace(")", "");
      sseDOM.innerText = date.toLocaleString();
      sseClientsDOM.innerText = activeClients;
      sseClientIDDOM.innerText = clientID;
    },
    onError: (err, { backoff }) => {
      sseClientsDOM.innerText = "N/A";

      let interval = null;
      interval = window.setInterval(() => {
        sseDOM.innerHTML = `SSE failed, attempting reconnect in <strong>${formatBackoff(
          backoff,
          0
        )}</strong>`;
        backoff -= 1000;
        if (backoff < 0) {
          sseDOM.innerHTML = `SSE failed, attempting reconnect in <strong>0s</strong>`;
          window.clearInterval(interval);
        }
      }, 1000);

      console.log(err);
    },
    initialBackoff: 1000,
    backoffStep: 1000,
  };

  const sseworker = new Worker("/static/js/sse.js");
  sseworker.onerror = (e) => {
    sseworker.terminate();
  };

  sseworker.onmessage = (e) => {
    if (e?.data?.error) {
      config.onError("SSE failed", e?.data);
    } else {
      config.onMessage(e?.data);
    }
  };

  sseworker.postMessage({
    url: config.url,
    initialBackoff: config.initialBackoff,
    backoffStep: config.backoffStep,
  });
};
webgo();
