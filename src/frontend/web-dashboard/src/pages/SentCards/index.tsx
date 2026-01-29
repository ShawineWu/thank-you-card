import React from "react";
import { SendCardForm } from "@/components/Cards/SendCardForm";
import { Send } from "lucide-react";

export const SentCards: React.FC = () => {
  const handleSendSuccess = () => {
    // Card sent successfully
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <div className="p-3 rounded-2xl bg-primary/10 text-primary">
          <Send size={24} />
        </div>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Sent Cards</h1>
          <p className="text-muted-foreground">
            The appreciation and gratitude you've shared with your team members.
          </p>
        </div>
      </div>

      <SendCardForm onSuccess={handleSendSuccess} />
    </div>
  );
};
