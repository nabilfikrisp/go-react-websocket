import { ChatRoom } from "./chat-room";
import { useChatSocket } from "./use-chat-socket";

type Props = {
  name: string;
};

export function ChatShell({ name }: Props) {
  const chat = useChatSocket(name);

  return <ChatRoom {...chat} />;
}
