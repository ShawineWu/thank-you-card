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
import { getValueBadgeClasses } from "@/constants/valueColors";

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
      <DialogContent className="sm:max-w-[500px] bg-white border-gray-200 text-gray-900">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-xl text-gray-800">
            <Heart className="h-5 w-5 text-indigo-500" fill="currentColor" />
            Recognition Detail
          </DialogTitle>
          <DialogDescription className="text-gray-500">
            A token of gratitude shared within the team.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6 pt-4">
          <div className="flex items-start justify-between">
            <div className="space-y-4">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-full bg-indigo-50 text-indigo-500">
                  <User size={18} />
                </div>
                <div>
                  <p className="text-xs font-medium text-gray-400 uppercase tracking-wider">
                    Sender
                  </p>
                  <div className="flex items-center gap-2 mt-1">
                    <Avatar className="h-6 w-6 border border-gray-200">
                      <AvatarFallback className="text-[10px] bg-indigo-50 text-indigo-600">
                        {getInitials(card.senderId)}
                      </AvatarFallback>
                    </Avatar>
                    <span className="font-semibold text-gray-800">{card.senderName || card.senderId}</span>
                  </div>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <div className="p-2 rounded-full bg-indigo-50 text-indigo-500">
                  <Calendar size={18} />
                </div>
                <div>
                  <p className="text-xs font-medium text-gray-400 uppercase tracking-wider">
                    Date Sent
                  </p>
                  <p className="mt-1 text-gray-700">
                    {format(new Date(card.createdAt), "MMMM do, yyyy")}
                  </p>
                </div>
              </div>
            </div>

            <div className="flex items-center justify-center p-3 rounded-2xl bg-gray-50 border border-gray-200">
              <div className="text-center">
                <p className="text-[10px] font-bold text-gray-400 uppercase">
                  Card ID
                </p>
                <p className="text-xs font-mono text-gray-600 mt-1">
                  #{card.id.substring(0, 8)}
                </p>
              </div>
            </div>
          </div>

          <div className="space-y-2">
            <div className="flex items-center gap-2 text-sm font-semibold text-gray-700">
              <Users size={16} className="text-indigo-500" />
              Recipients
            </div>
            <div className="flex flex-col gap-2">
              {card.recipients.map((recipient) => (
                <div key={recipient.id} className="flex items-center gap-2">
                  <Avatar className="h-6 w-6 border border-gray-200">
                    <AvatarFallback className="text-[10px] bg-pink-50 text-pink-600">
                      {getInitials(recipient.name || recipient.id)}
                    </AvatarFallback>
                  </Avatar>
                  <span className="text-sm text-gray-700">{recipient.name}</span>
                </div>
              ))}
            </div>
          </div>

          <div className="space-y-3">
            <div className="flex items-center gap-2 text-sm font-semibold text-gray-700">
              <Quote size={16} className="text-indigo-500" />
              Recognition Reason
            </div>
            <div className="p-4 rounded-xl bg-gray-50 border border-gray-200 italic text-gray-600 leading-relaxed">
              "{card.recognitionReason}"
            </div>
          </div>

          <div className="space-y-3">
            <p className="text-sm font-semibold text-gray-700">
              Associated Values
            </p>
            <div className="flex flex-wrap gap-2">
              {card.selectedValues.map((val) => (
                <Badge
                  key={val.id}
                  className={`border ${getValueBadgeClasses(val.name)}`}
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
