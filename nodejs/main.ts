import { WebSocket } from "ws";

const accessToken = process.env.TRAQ_BOT_ACCESS_TOKEN;
if (!accessToken) throw "TRAQ_BOT_ACCESS_TOKEN is not set";

const handleMessage = async (message: string) => {
  const payload = JSON.parse(message) as
    | {
        type: "MESSAGE_CREATED";
        reqId: string;
        body: {
          message: {
            user: {
              name: string;
              bot: boolean;
            };
            channelId: string;
            plainText: string;
          };
        };
      }
    | {
        type: string; // TODO: except "MESSAGE_CREATED"
        reqId: string;
      };

  if (payload.type !== "MESSAGE_CREATED" || !("body" in payload)) {
    console.log(`unsupported events(${payload.reqId}): ${payload.type}`);
    return;
  }

  const { user, channelId, plainText } = payload.body.message;

  if (user.bot) {
    console.log(`bot message(${payload.reqId})`);
    return;
  }

  const args = plainText.split(" ");
  if (args.length !== 2 || !args[0].startsWith("@") || args[1] !== "nodejs") {
    console.log(`invalid args(${payload.reqId}): ${plainText}`);
    return;
  }

  const stamp = ":node_js:";
  const content = `@${user.name} ${stamp}`;
  await postMessage(channelId, content);
};

const postMessage = async (channelId: string, content: string) => {
  const body = JSON.stringify({
    content,
    embed: true,
  });

  await fetch(`https://q.trap.jp/api/v3/channels/${channelId}/messages`, {
    method: "POST",
    body,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${accessToken}`,
    },
  });
};

if (require.main) {
  const ws = new WebSocket("wss://q.trap.jp/api/v3/bots/ws", {
    headers: {
      authorization: `Bearer ${accessToken}`,
    },
  });

  console.log("connected");

  ws.onmessage = (e) => handleMessage(e.data.toString());
}
