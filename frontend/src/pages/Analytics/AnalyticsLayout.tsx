import { Outlet, Link, useLocation } from 'react-router-dom'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { BarChart3, Users, Building2, PieChart } from 'lucide-react'

const AnalyticsLayout = () => {
  const location = useLocation()

  const tabs = [
    { path: '/analytics/dashboard', label: 'Dashboard', icon: BarChart3 },
    { path: '/analytics/recognizers', label: 'Recognizers', icon: Users },
    { path: '/analytics/teams', label: 'Teams', icon: Building2 },
    { path: '/analytics/values', label: 'Values', icon: PieChart },
  ]

  const currentTab = tabs.find(tab => location.pathname === tab.path)?.path || tabs[0].path

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">HR Analytics</h1>
        <p className="text-muted-foreground">
          Comprehensive insights into recognition patterns across the company
        </p>
      </div>

      <Tabs value={currentTab}>
        <TabsList className="grid w-full grid-cols-4">
          {tabs.map((tab) => {
            const Icon = tab.icon
            return (
              <Link key={tab.path} to={tab.path}>
                <TabsTrigger value={tab.path} className="flex items-center space-x-2 w-full">
                  <Icon className="h-4 w-4" />
                  <span className="hidden sm:inline">{tab.label}</span>
                </TabsTrigger>
              </Link>
            )
          })}
        </TabsList>
      </Tabs>

      <div className="mt-6">
        <Outlet />
      </div>
    </div>
  )
}

export default AnalyticsLayout
