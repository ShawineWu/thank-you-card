import React, { useState, useEffect } from "react";
import {
  Drawer,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
} from "@/components/ui/drawer";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { api, type Employee, type ValueResponse } from "@/services/api";
import { useAuthStore } from "@/store/authStore";
import { X, Search, Loader2, Send, Sparkles } from "lucide-react";
import { useToast } from "@/hooks/use-toast";

interface SendCardDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

interface Recipient {
  id: string;
  name: string;
}

export const SendCardDrawer: React.FC<SendCardDrawerProps> = ({
  open,
  onOpenChange,
  onSuccess,
}) => {
  const { toast } = useToast();
  const { user } = useAuthStore();
  const senderName =
    user?.profile?.name || user?.profile?.preferred_username || "User";

  const [recipients, setRecipients] = useState<Recipient[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<Employee[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [showSearchResults, setShowSearchResults] = useState(false);
  const [keyInfo, setKeyInfo] = useState("");
  const [recognitionReason, setRecognitionReason] = useState("");
  const [selectedValueIds, setSelectedValueIds] = useState<string[]>([]);
  const [values, setValues] = useState<ValueResponse[]>([]);
  const [loadingValues, setLoadingValues] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isGenerating, setIsGenerating] = useState(false);
  const [errors, setErrors] = useState<{
    recipients?: string;
    recognitionReason?: string;
    valueIds?: string;
  }>({});

  // Load company values
  useEffect(() => {
    if (open) {
      const loadValues = async () => {
        setLoadingValues(true);
        try {
          const response = await api.getValues();
          setValues(response.data.data);
        } catch (error: any) {
          toast({
            title: "Error",
            description: "Failed to load company values",
            variant: "destructive",
          });
        } finally {
          setLoadingValues(false);
        }
      };
      loadValues();
    }
  }, [open, toast]);

  // Debounced employee search
  useEffect(() => {
    if (!open) return;

    const searchEmployees = async () => {
      if (searchQuery.trim().length < 2) {
        setSearchResults([]);
        setShowSearchResults(false);
        return;
      }

      setIsSearching(true);
      try {
        const response = await api.searchEmployees(searchQuery.trim());
        setSearchResults(response.data.data);
        setShowSearchResults(true);
      } catch (error: any) {
        console.error("Failed to search employees", error);
        setSearchResults([]);
      } finally {
        setIsSearching(false);
      }
    };

    const timeoutId = setTimeout(searchEmployees, 300);
    return () => clearTimeout(timeoutId);
  }, [searchQuery, open]);

  const handleSelectRecipient = (employee: Employee) => {
    if (!recipients.find((r) => r.id === employee.id)) {
      setRecipients([...recipients, { id: employee.id, name: employee.name }]);
      setSearchQuery("");
      setShowSearchResults(false);
      if (errors.recipients) {
        setErrors({ ...errors, recipients: undefined });
      }
    }
  };

  const handleRemoveRecipient = (id: string) => {
    setRecipients(recipients.filter((r) => r.id !== id));
  };

  const handleToggleValue = (valueId: string) => {
    if (selectedValueIds.includes(valueId)) {
      setSelectedValueIds(selectedValueIds.filter((id) => id !== valueId));
    } else {
      if (selectedValueIds.length < 3) {
        setSelectedValueIds([...selectedValueIds, valueId]);
      }
    }
    if (errors.valueIds) {
      setErrors({ ...errors, valueIds: undefined });
    }
  };

  const handleGenerateReason = async () => {
    if (!keyInfo.trim()) {
      toast({
        title: "Key information required",
        description:
          "Please enter key information to generate recognition reason.",
        variant: "destructive",
      });
      return;
    }
    if (selectedValueIds.length === 0) {
      toast({
        title: "Company values required",
        description: "Please select at least one company value.",
        variant: "destructive",
      });
      return;
    }
    if (recipients.length === 0) {
      toast({
        title: "Recipients required",
        description: "Please select at least one recipient.",
        variant: "destructive",
      });
      return;
    }

    setIsGenerating(true);
    try {
      const selectedValues = values.filter((v) =>
        selectedValueIds.includes(v.id)
      );
      const valueNames = selectedValues.map((v) => v.name);
      const recipientNames = recipients.map((r) => r.name);
      const res = await api.generateRecognitionReason({
        keyInfo: keyInfo.trim(),
        valueNames,
        senderName,
        recipientNames,
      });
      const generated = res.data.data?.recognitionReason ?? "";
      setRecognitionReason(generated);
      if (errors.recognitionReason) {
        setErrors({ ...errors, recognitionReason: undefined });
      }
      toast({
        title: "Generated",
        description:
          "Recognition reason has been generated. You can edit it below.",
      });
    } catch (error: any) {
      toast({
        title: "Generation failed",
        description:
          error.response?.data?.error?.message ||
          "Failed to generate recognition reason.",
        variant: "destructive",
      });
    } finally {
      setIsGenerating(false);
    }
  };

  const validate = (): boolean => {
    const newErrors: typeof errors = {};

    if (recipients.length === 0) {
      newErrors.recipients = "Please select at least one recipient";
    }

    if (!recognitionReason.trim()) {
      newErrors.recognitionReason = "Please provide a recognition reason";
    } else if (recognitionReason.trim().length < 30) {
      newErrors.recognitionReason = "Reason must be at least 30 characters";
    } else if (recognitionReason.length > 1000) {
      newErrors.recognitionReason = "Reason must not exceed 1000 characters";
    }

    if (selectedValueIds.length === 0) {
      newErrors.valueIds = "Please select at least one company value";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) {
      return;
    }

    setIsSubmitting(true);
    try {
      await api.createCard({
        recipients: recipients,
        recognitionReason: recognitionReason.trim(),
        valueIds: selectedValueIds,
      });

      toast({
        title: "Success",
        description: "Recognition card sent successfully! 🎉",
      });

      // Reset form
      setRecipients([]);
      setRecognitionReason("");
      setSelectedValueIds([]);
      setKeyInfo("");
      setErrors({});
      setSearchQuery("");
      setShowSearchResults(false);

      onOpenChange(false);
      onSuccess?.();
    } catch (error: any) {
      toast({
        title: "Error",
        description:
          error.response?.data?.error?.message ||
          "Failed to send recognition card",
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    if (!isSubmitting) {
      setRecipients([]);
      setRecognitionReason("");
      setSelectedValueIds([]);
      setKeyInfo("");
      setErrors({});
      setSearchQuery("");
      setShowSearchResults(false);
      onOpenChange(false);
    }
  };

  return (
    <Drawer open={open} onOpenChange={handleClose}>
      <DrawerContent>
        <DrawerHeader>
          <DrawerTitle>Send Recognition Card</DrawerTitle>
          <DrawerDescription>
            Appreciate your colleagues and celebrate our company values
          </DrawerDescription>
        </DrawerHeader>

        <form onSubmit={handleSubmit} className="flex flex-col flex-1">
          <div className="flex-1 px-6 space-y-6 overflow-y-auto">
            {/* Recipients */}
            <div className="space-y-2">
              <label className="text-sm font-medium">
                Recipients <span className="text-red-500">*</span>
              </label>
              <div className="relative">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    type="text"
                    placeholder="Search for colleagues..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    onFocus={() => {
                      if (searchResults.length > 0) {
                        setShowSearchResults(true);
                      }
                    }}
                    className="pl-9"
                  />
                </div>

                {showSearchResults && (
                  <div className="absolute z-10 w-full mt-1 bg-popover border rounded-md shadow-lg max-h-60 overflow-y-auto">
                    {isSearching ? (
                      <div className="p-4 text-center text-sm text-muted-foreground">
                        <Loader2 className="h-4 w-4 animate-spin inline mr-2" />
                        Searching...
                      </div>
                    ) : searchResults.length > 0 ? (
                      searchResults.map((emp) => (
                        <div
                          key={emp.id}
                          onClick={() => handleSelectRecipient(emp)}
                          className="p-3 hover:bg-accent cursor-pointer border-b last:border-b-0"
                        >
                          <div className="font-medium">{emp.name}</div>
                          <div className="text-sm text-muted-foreground">
                            {emp.email}
                            {emp.department && ` · ${emp.department}`}
                          </div>
                        </div>
                      ))
                    ) : searchQuery.trim().length >= 2 ? (
                      <div className="p-4 text-center text-sm text-muted-foreground">
                        No employees found
                      </div>
                    ) : null}
                  </div>
                )}
              </div>

              {recipients.length > 0 && (
                <div className="flex flex-wrap gap-2 mt-2">
                  {recipients.map((r) => (
                    <Badge
                      key={r.id}
                      variant="secondary"
                      className="px-3 py-1 flex items-center gap-2"
                    >
                      {r.name}
                      <button
                        type="button"
                        onClick={() => handleRemoveRecipient(r.id)}
                        className="ml-1 hover:text-destructive"
                      >
                        <X className="h-3 w-3" />
                      </button>
                    </Badge>
                  ))}
                </div>
              )}

              {errors.recipients && (
                <p className="text-sm text-destructive">{errors.recipients}</p>
              )}
            </div>

            {/* Company Values */}
            <div className="space-y-2">
              <label className="text-sm font-medium">
                Company Values <span className="text-red-500">*</span>
              </label>
              {loadingValues ? (
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Loading values...
                </div>
              ) : (
                <div className="flex flex-wrap gap-2">
                  {values.map((value) => {
                    const isSelected = selectedValueIds.includes(value.id);
                    const isDisabled =
                      !isSelected && selectedValueIds.length >= 3;

                    return (
                      <button
                        key={value.id}
                        type="button"
                        onClick={() => handleToggleValue(value.id)}
                        disabled={isDisabled}
                        className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${
                          isSelected
                            ? "bg-primary text-primary-foreground"
                            : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
                        } ${
                          isDisabled
                            ? "opacity-50 cursor-not-allowed"
                            : "cursor-pointer"
                        }`}
                        title={value.description}
                      >
                        {value.name}
                      </button>
                    );
                  })}
                </div>
              )}
              <div className="flex justify-between">
                <p className="text-xs text-muted-foreground">
                  Select 1-3 values ({selectedValueIds.length}/3 selected)
                </p>
                {errors.valueIds && (
                  <p className="text-xs text-destructive">{errors.valueIds}</p>
                )}
              </div>
            </div>

            {/* Key information */}
            <div className="space-y-2">
              <label className="text-sm font-medium">Key information</label>
              <textarea
                value={keyInfo}
                onChange={(e) => setKeyInfo(e.target.value)}
                placeholder="Enter key points (e.g. what they did, project name, impact). Used to generate Recognition Reason."
                rows={3}
                className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 resize-none"
              />
              <div className="flex justify-end">
                <Button
                  type="button"
                  variant="secondary"
                  onClick={handleGenerateReason}
                  disabled={
                    isGenerating ||
                    !keyInfo.trim() ||
                    selectedValueIds.length === 0 ||
                    recipients.length === 0
                  }
                >
                  {isGenerating ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Generating...
                    </>
                  ) : (
                    <>
                      <Sparkles className="mr-2 h-4 w-4" />
                      Generate
                    </>
                  )}
                </Button>
              </div>
            </div>

            {/* Recognition Reason */}
            <div className="space-y-2">
              <label className="text-sm font-medium">
                Recognition Reason <span className="text-red-500">*</span>
              </label>
              <textarea
                value={recognitionReason}
                onChange={(e) => {
                  setRecognitionReason(e.target.value);
                  if (errors.recognitionReason) {
                    setErrors({ ...errors, recognitionReason: undefined });
                  }
                }}
                placeholder="Describe why you're recognizing this person, or use Generate above..."
                rows={5}
                className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 resize-none"
              />
              <div className="flex justify-between">
                <p className="text-xs text-muted-foreground">
                  {recognitionReason.length}/1000 characters
                </p>
                {errors.recognitionReason && (
                  <p className="text-xs text-destructive">
                    {errors.recognitionReason}
                  </p>
                )}
              </div>
            </div>
          </div>

          <DrawerFooter>
            <Button
              type="button"
              variant="outline"
              onClick={handleClose}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Sending...
                </>
              ) : (
                <>
                  <Send className="mr-2 h-4 w-4" />
                  Send Recognition
                </>
              )}
            </Button>
          </DrawerFooter>
        </form>
      </DrawerContent>
    </Drawer>
  );
};
