// chat-room.tsx
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import type { ChatMessage } from "./use-chat-socket";

type Props = {
  status: "connecting" | "open" | "closed" | "error";
  messages: ChatMessage[];
  send: (payload: { content: string }) => void;
};

export function ChatRoom({ status, messages, send }: Props) {
  const [text, setText] = useState("");
  const disabled = status !== "open";

  const handleSend = () => {
    if (disabled) return;

    const trimmed = text.trim();
    if (!trimmed) return;

    send({ content: trimmed });
    setText("");
  };

  return (
    <div className="flex h-screen flex-col p-4 gap-4">
      <div className="text-sm text-muted-foreground">Status: {status}</div>

      <ul className="flex-1 overflow-y-auto space-y-1">
        {messages.map((m, i) => (
          <li key={i}>
            <strong>{m.sender}:</strong> {m.content}
          </li>
        ))}
      </ul>

      <div className="flex gap-2">
        <Input
          value={text}
          onChange={(e) => setText(e.target.value)}
          disabled={disabled}
        />
        <Button onClick={handleSend} disabled={disabled}>
          Send
        </Button>
      </div>
    </div>
  );
}
