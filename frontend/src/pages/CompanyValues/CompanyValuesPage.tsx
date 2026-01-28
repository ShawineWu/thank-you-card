import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { companyValuesApi } from '@/services/companyValues'
import { CompanyValueDetail } from '@/types/card'
import ValueDetailDialog from '@/components/CompanyValues/ValueDetailDialog'
import { Info } from 'lucide-react'

const CompanyValuesPage = () => {
  const [selectedValue, setSelectedValue] = useState<CompanyValueDetail | null>(null)
  const [dialogOpen, setDialogOpen] = useState(false)

  const { data: allData, isLoading: isLoadingAll } = useQuery({
    queryKey: ['company-values', 'all'],
    queryFn: () => companyValuesApi.getAll(),
  })

  const { data: valuesData, isLoading: isLoadingValues } = useQuery({
    queryKey: ['company-values', 'VALUE'],
    queryFn: () => companyValuesApi.getByType('VALUE'),
  })

  const { data: credosData, isLoading: isLoadingCredos } = useQuery({
    queryKey: ['company-values', 'CREDO'],
    queryFn: () => companyValuesApi.getByType('CREDO'),
  })

  const handleViewDetail = (value: CompanyValueDetail) => {
    setSelectedValue(value)
    setDialogOpen(true)
  }

  const renderValueCard = (value: CompanyValueDetail) => (
    <Card key={value.id} className="hover:shadow-md transition-shadow">
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <CardTitle className="text-lg mb-2">{value.name}</CardTitle>
            <CardDescription className="text-xs text-muted-foreground">
              {value.code}
            </CardDescription>
          </div>
          <Badge variant={value.type === 'VALUE' ? 'default' : 'success'} className="ml-2">
            {value.type}
          </Badge>
        </div>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-muted-foreground line-clamp-3 mb-4">
          {value.description}
        </p>
        <Button
          variant="outline"
          size="sm"
          onClick={() => handleViewDetail(value)}
          className="w-full"
        >
          <Info className="h-4 w-4 mr-2" />
          View Details
        </Button>
      </CardContent>
    </Card>
  )

  if (isLoadingAll || isLoadingValues || isLoadingCredos) {
    return (
      <div className="text-center py-12">
        <p className="text-muted-foreground">Loading company values...</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold mb-2">Company Values & Credos</h1>
        <p className="text-muted-foreground">
          Explore our company values and credos that guide our work and culture.
        </p>
      </div>

      <Tabs defaultValue="all" className="w-full">
        <TabsList className="grid w-full max-w-md grid-cols-3">
          <TabsTrigger value="all">
            All ({allData?.total || 0})
          </TabsTrigger>
          <TabsTrigger value="values">
            Values ({valuesData?.total || 0})
          </TabsTrigger>
          <TabsTrigger value="credos">
            Credos ({credosData?.total || 0})
          </TabsTrigger>
        </TabsList>

        <TabsContent value="all" className="mt-6">
          {allData && allData.items.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {allData.items.map(renderValueCard)}
            </div>
          ) : (
            <div className="text-center py-12">
              <p className="text-muted-foreground">No company values found.</p>
            </div>
          )}
        </TabsContent>

        <TabsContent value="values" className="mt-6">
          {valuesData && valuesData.items.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {valuesData.items.map(renderValueCard)}
            </div>
          ) : (
            <div className="text-center py-12">
              <p className="text-muted-foreground">No values found.</p>
            </div>
          )}
        </TabsContent>

        <TabsContent value="credos" className="mt-6">
          {credosData && credosData.items.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {credosData.items.map(renderValueCard)}
            </div>
          ) : (
            <div className="text-center py-12">
              <p className="text-muted-foreground">No credos found.</p>
            </div>
          )}
        </TabsContent>
      </Tabs>

      <ValueDetailDialog
        value={selectedValue}
        open={dialogOpen}
        onOpenChange={setDialogOpen}
      />
    </div>
  )
}

export default CompanyValuesPage
