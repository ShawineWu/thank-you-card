import React from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import type { Card as CardType } from "@/services/api";
import { format } from "date-fns";
import { Calendar, User, Users, Quote, Heart } from "lucide-react";

interface CardDetailProps {
  card: CardType | null;
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
}

export const CardDetail: React.FC<CardDetailProps> = ({
  card,
  isOpen,
  onOpenChange,
}) => {
  if (!card) return null;

  const getInitials = (id: string) => id.substring(0, 2).toUpperCase();

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px] bg-slate-950 border-slate-800 text-slate-100">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-xl">
            <Heart className="h-5 w-5 text-primary" fill="currentColor" />
            Recognition Detail
          </DialogTitle>
          <DialogDescription className="text-slate-400">
            A token of gratitude shared within the team.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6 pt-4">
          <div className="flex items-start justify-between">
            <div className="space-y-4">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-full bg-primary/10 text-primary">
                  <User size={18} />
                </div>
                <div>
                  <p className="text-xs font-medium text-slate-500 uppercase tracking-wider">
                    Sender
                  </p>
                  <div className="flex items-center gap-2 mt-1">
                    <Avatar className="h-6 w-6">
                      <AvatarFallback className="text-[10px] bg-primary/20">
                        {getInitials(card.senderId)}
                      </AvatarFallback>
                    </Avatar>
                    <span className="font-semibold">{card.senderName || card.senderId}</span>
                  </div>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <div className="p-2 rounded-full bg-primary/10 text-primary">
                  <Calendar size={18} />
                </div>
                <div>
                  <p className="text-xs font-medium text-slate-500 uppercase tracking-wider">
                    Date Sent
                  </p>
                  <p className="mt-1">
                    {format(new Date(card.createdAt), "MMMM do, yyyy")}
                  </p>
                </div>
              </div>
            </div>

            <div className="flex items-center justify-center p-3 rounded-2xl bg-slate-900 border border-slate-800">
              <div className="text-center">
                <p className="text-[10px] font-bold text-slate-500 uppercase">
                  Card ID
                </p>
                <p className="text-xs font-mono text-slate-300 mt-1">
                  #{card.id.substring(0, 8)}
                </p>
              </div>
            </div>
          </div>

          <div className="space-y-2">
            <div className="flex items-center gap-2 text-sm font-semibold text-slate-300">
              <Users size={16} className="text-primary" />
              Recipients
            </div>
            <div className="flex flex-wrap gap-2">
              {card.recipients.map((recipient) => (
                <Badge
                  key={recipient.id}
                  variant="secondary"
                  className="bg-slate-900 text-slate-300 hover:bg-slate-800 border-slate-700"
                >
                  {recipient.name}
                </Badge>
              ))}
            </div>
          </div>

          <div className="space-y-3">
            <div className="flex items-center gap-2 text-sm font-semibold text-slate-300">
              <Quote size={16} className="text-primary" />
              Recognition Reason
            </div>
            <div className="p-4 rounded-xl bg-slate-900/50 border border-slate-800 italic text-slate-300 leading-relaxed shadow-inner">
              "{card.recognitionReason}"
            </div>
          </div>

          <div className="space-y-3">
            <p className="text-sm font-semibold text-slate-300">
              Associated Values
            </p>
            <div className="flex flex-wrap gap-2">
              {card.selectedValues.map((val) => (
                <Badge
                  key={val.id}
                  className="bg-primary/20 text-primary border-primary/30 hover:bg-primary/30"
                >
                  {val.name}
                </Badge>
              ))}
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};
