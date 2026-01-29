import React, { useEffect, useState } from "react";
import { api, type ValueResponse } from "@/services/api";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { Award, BookOpen, ChevronRight } from "lucide-react";

const VALUE_TYPE = "Value";
const CREDO_TYPE = "Credo";

export const CompanyValues: React.FC = () => {
  const [values, setValues] = useState<ValueResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [detailValue, setDetailValue] = useState<ValueResponse | null>(null);
  const [detailOpen, setDetailOpen] = useState(false);

  useEffect(() => {
    const fetchValues = async () => {
      try {
        setLoading(true);
        const res = await api.getValues();
        setValues(res.data.data ?? []);
      } catch (error) {
        console.error("Failed to fetch company values:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchValues();
  }, []);

  const valueItems = values.filter((v) => v.type === VALUE_TYPE);
  const credoItems = values.filter((v) => v.type === CREDO_TYPE);

  const openDetails = (v: ValueResponse) => {
    setDetailValue(v);
    setDetailOpen(true);
  };

  const ValueCard: React.FC<{
    item: ValueResponse;
    onViewDetails: () => void;
  }> = ({ item, onViewDetails }) => (
    <Card className="overflow-hidden transition-shadow hover:shadow-md">
      <CardHeader className="pb-2">
        <CardTitle className="text-lg">{item.name}</CardTitle>
        <CardDescription className="line-clamp-2">{item.description}</CardDescription>
      </CardHeader>
      <CardContent className="pt-0">
        <Button variant="outline" size="sm" onClick={onViewDetails}>
          View details
          <ChevronRight className="ml-1 h-4 w-4" />
        </Button>
      </CardContent>
    </Card>
  );

  const SectionSkeleton: React.FC<{ count?: number }> = ({ count = 4 }) => (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {Array.from({ length: count }).map((_, i) => (
        <Skeleton key={i} className="h-[140px] rounded-lg" />
      ))}
    </div>
  );

  return (
    <div className="space-y-8">
      <div className="flex items-center gap-3">
        <div className="p-3 rounded-2xl bg-primary/10 text-primary">
          <Award size={24} />
        </div>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Company Values</h1>
          <p className="text-muted-foreground">
            Explore our values and credos that guide how we work and recognize each other.
          </p>
        </div>
      </div>

      {loading ? (
        <>
          <div>
            <h2 className="flex items-center gap-2 text-xl font-semibold mb-4">
              <BookOpen className="h-5 w-5" /> Values
            </h2>
            <SectionSkeleton />
          </div>
          <div>
            <h2 className="flex items-center gap-2 text-xl font-semibold mb-4">
              <Award className="h-5 w-5" /> Credo
            </h2>
            <SectionSkeleton count={3} />
          </div>
        </>
      ) : (
        <>
          <div>
            <h2 className="flex items-center gap-2 text-xl font-semibold mb-4">
              <BookOpen className="h-5 w-5" /> Values
            </h2>
            {valueItems.length === 0 ? (
              <p className="text-muted-foreground">No values defined.</p>
            ) : (
              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {valueItems.map((v) => (
                  <ValueCard
                    key={v.id}
                    item={v}
                    onViewDetails={() => openDetails(v)}
                  />
                ))}
              </div>
            )}
          </div>

          <div>
            <h2 className="flex items-center gap-2 text-xl font-semibold mb-4">
              <Award className="h-5 w-5" /> Credo
            </h2>
            {credoItems.length === 0 ? (
              <p className="text-muted-foreground">No credos defined.</p>
            ) : (
              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {credoItems.map((v) => (
                  <ValueCard
                    key={v.id}
                    item={v}
                    onViewDetails={() => openDetails(v)}
                  />
                ))}
              </div>
            )}
          </div>
        </>
      )}

      <Dialog open={detailOpen} onOpenChange={setDetailOpen}>
        <DialogContent className="max-w-md sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>{detailValue?.name}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <p className="text-sm font-medium text-muted-foreground mb-1">Description</p>
              <p className="text-sm">{detailValue?.description}</p>
            </div>
            {detailValue?.examples && detailValue.examples.length > 0 && (
              <div>
                <p className="text-sm font-medium text-muted-foreground mb-2">Examples</p>
                <ul className="list-disc list-inside space-y-1 text-sm">
                  {detailValue.examples.map((ex, i) => (
                    <li key={i}>{ex}</li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
};
