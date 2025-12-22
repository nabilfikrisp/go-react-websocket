// use-chat-socket.ts
import { useEffect, useRef, useState } from "react";

export type ChatMessage = {
  sender: string;
  content: string;
};

type Status = "connecting" | "open" | "closed" | "error";

export const CHAT_WS_ENDPOINT = "ws://localhost:3001/ws";

export function useChatSocket(name: string) {
  const socketRef = useRef<WebSocket | null>(null);
  const didUnmountRef = useRef(false);
  const [status, setStatus] = useState<Status>("connecting");
  const [messages, setMessages] = useState<ChatMessage[]>([]);

  useEffect(() => {
    console.log("useChatSocket MOUNT");

    const socket = new WebSocket(`${CHAT_WS_ENDPOINT}?name=${name}`);
    socketRef.current = socket;
    didUnmountRef.current = false;

    socket.onopen = () => {
      console.log("WS OPEN");

      if (didUnmountRef.current) return;
      setStatus("open");
    };

    socket.onmessage = (event) => {
      if (didUnmountRef.current) return;
      if (socket.readyState !== WebSocket.OPEN) return;

      const data = JSON.parse(event.data);
      setMessages((prev) => [...prev, data]);
    };

    socket.onerror = () => {
      if (didUnmountRef.current) return;
      setStatus((prev) => (prev === "open" ? "error" : prev));
    };

    socket.onclose = () => {
      console.log("WS CLOSE");

      if (didUnmountRef.current) return;
      setStatus((prev) => (prev === "error" ? prev : "closed"));
    };

    return () => {
      console.log("useChatSocket CLEANUP");

      didUnmountRef.current = true;

      // remove handlers first
      socket.onopen = null;
      socket.onmessage = null;
      socket.onerror = null;
      socket.onclose = null;

      // close exactly once
      if (
        socket.readyState === WebSocket.OPEN ||
        socket.readyState === WebSocket.CONNECTING
      ) {
        socket.close();
      }

      socketRef.current = null;
    };
  }, [name]);

  const send = (payload: { content: string }) => {
    if (socketRef.current?.readyState !== WebSocket.OPEN) return;
    socketRef.current.send(JSON.stringify(payload));
  };

  return {
    status,
    messages,
    send,
  };
}
