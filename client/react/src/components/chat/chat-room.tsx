// chat-room.tsx
import { useState, useRef, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { HugeiconsIcon } from "@hugeicons/react";
import { SentIcon } from "@hugeicons/core-free-icons";
import type { ChatMessage } from "./use-chat-socket";

type Props = {
  status: "connecting" | "open" | "closed" | "error";
  messages: ChatMessage[];
  send: (payload: { content: string }) => void;
  currentUser?: string;
};

export function ChatRoom({ status, messages, send, currentUser }: Props) {
  const [text, setText] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const disabled = status !== "open";

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSend = () => {
    if (disabled) return;

    const trimmed = text.trim();
    if (!trimmed) return;

    send({ content: trimmed });
    setText("");
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const getStatusVariant = () => {
    switch (status) {
      case "open":
        return "default";
      case "connecting":
        return "secondary";
      case "error":
        return "destructive";
      default:
        return "outline";
    }
  };

  return (
    <div className="flex h-screen flex-col gap-4 p-4">
      <Card className="flex-1 flex flex-col">
        <CardHeader className="border-b">
          <div className="flex items-center justify-between">
            <CardTitle>Chat Room</CardTitle>
            <Badge variant={getStatusVariant()}>
              {status.charAt(0).toUpperCase() + status.slice(1)}
            </Badge>
          </div>
        </CardHeader>
        <CardContent className="flex-1 flex flex-col gap-4 overflow-hidden">
          <div className="flex-1 overflow-y-auto space-y-3">
            {messages.map((m, i) => {
              const isOwn = m.sender === currentUser;
              return (
                <div
                  key={i}
                  className={`flex ${isOwn ? "justify-end" : "justify-start"}`}
                >
                  <div
                    className={`flex max-w-[75%] flex-col gap-1 ${
                      isOwn ? "items-end" : "items-start"
                    }`}
                  >
                    {!isOwn && (
                      <span className="text-xs font-medium text-muted-foreground px-1">
                        {m.sender}
                      </span>
                    )}
                    <Card
                      className={`px-3 py-2 ${
                        isOwn
                          ? "bg-primary text-primary-foreground border-primary"
                          : "border ring-0"
                      }`}
                    >
                      <p className="text-sm whitespace-pre-wrap wrap-break-word">
                        {m.content}
                      </p>
                    </Card>
                  </div>
                </div>
              );
            })}
            <div ref={messagesEndRef} />
          </div>

          <div className="flex gap-2">
            <Input
              value={text}
              onChange={(e) => setText(e.target.value)}
              onKeyDown={handleKeyPress}
              disabled={disabled}
              placeholder={disabled ? "Connecting..." : "Type a message..."}
              className="flex-1"
            />
            <Button
              onClick={handleSend}
              disabled={disabled || !text.trim()}
              size="icon"
            >
              <HugeiconsIcon icon={SentIcon} strokeWidth={2} />
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
