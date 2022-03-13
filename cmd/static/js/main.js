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
      source.onopen = () => {
        // reset backoff to initial, so further failures will again start with initial backoff
        // instead of previous duration
        backoff = initialBackoff;
      };

      source.onmessage = (event) => {
        onMessage && onMessage(event);
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
        onError && onError(err);
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

  sse(`/sse/${clientID}`, {
    onMessage: (event) => {
      const parts = event.data?.split("(");
      const date = new Date(parts[0]);
      const activeClients = parts[1].replace(")", "");
      sseDOM.innerText = date.toLocaleString();
      sseClientsDOM.innerText = activeClients;
      sseClientIDDOM.innerText = clientID;
    },
    onError: (err) => {
      console.log(err);
      sseDOM.innerText = `SSE error, restarting`;
    },
    backoffStep: 150,
  })();
};
webgo();
