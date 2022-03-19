const webgo = async () => {
  const sse = (url, config = {}) => {
    const {
      onMessage,
      onError,
      initialBackoff = 10, // milliseconds
      maxBackoff = 15 * 1000, // 15 seconds
      backoffStep = 50, // milliseconds
    } = config;

    let backoff = initialBackoff,
      sseRetryTimeout = null;

    const start = () => {
      const source = new EventSource(url);
      const configState = { initialBackoff, maxBackoff, backoffStep, backoff };

      source.onopen = () => {
        clearTimeout(sseRetryTimeout);
        // reset backoff to initial, so further failures will again start with initial backoff
        // instead of previous duration
        backoff = initialBackoff;
        configState.backoff = backoff
      };

      source.onmessage = (event, configState) => {
        onMessage && onMessage(event, configState);
      };

      source.onerror = (err) => {
        source.close();
        clearTimeout(sseRetryTimeout);
        // reattempt connecting with *linear* backoff
        sseRetryTimeout = window.setTimeout(() => {
          start(url, onMessage);
          if (backoff < maxBackoff) {
            backoff += backoffStep;
            if (backoff > maxBackoff) {
              backoff = maxBackoff;
            }
          }
        }, backoff);
        onError && onError(err, configState);
      };
    };
    return start;
  };

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

  sse(`/sse/${clientID}`, {
    onMessage: (event) => {
      const parts = event.data?.split("(");
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
        if (backoff <  0) {
          sseDOM.innerHTML = `SSE failed, attempting reconnect in <strong>0s</strong>`;
          window.clearInterval(interval);
        }
      }, 1000);

      console.log(err);
    },
    initialBackoff: 1000,
    backoffStep: 1000,
  })();
};
webgo();
