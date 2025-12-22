import { useState } from "react";
import { NameGate } from "./components/chat/name-gate";
import { ChatShell } from "./components/chat/chat-shell";

export default function App() {
  const [name, setName] = useState<string | null>(null);

  if (name === null) {
    return <NameGate onJoin={setName} />;
  }

  return <ChatShell name={name} />;
}
