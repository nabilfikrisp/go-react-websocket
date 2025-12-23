import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

type Props = {
  onJoin: (name: string) => void;
};

export function NameGate({ onJoin }: Props) {
  const [name, setName] = useState("");

  const handleJoin = () => {
    const trimmed = name.trim();
    if (!trimmed) return;
    onJoin(trimmed);
  };

  return (
    <div className="flex h-screen items-center justify-center">
      <div className="flex flex-col gap-4">
        <Input
          placeholder="Enter your name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <Button onClick={handleJoin} className="ml-auto">
          Join
        </Button>
      </div>
    </div>
  );
}
