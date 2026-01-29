import { Outlet, Link, useLocation, useNavigate } from 'react-router-dom'
import { Button } from './ui/button'
import { Heart, Home, Plus, LayoutDashboard, Settings, BookOpen, LogOut } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'

const Layout = () => {
  const location = useLocation()
  const { isAuthenticated, isHROrAdmin, user, logout } = useAuth()
  const navigate = useNavigate()

  const navItems = [
    { path: '/', label: 'Overview', icon: LayoutDashboard },
    { path: '/feed', label: 'Feed', icon: Home },
    { path: '/create', label: 'Create Card', icon: Plus },
    { path: '/company-values', label: 'Values', icon: BookOpen },
  ]

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const isActive = (path: string) => {
    if (path === '/') {
      return location.pathname === '/' || location.pathname === '/overview'
    }
    return location.pathname.startsWith(path)
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b bg-card">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <Link to="/" className="flex items-center space-x-2">
              <Heart className="h-6 w-6 text-primary" />
              <span className="text-xl font-bold">Thank You Card</span>
            </Link>
            <nav className="hidden md:flex items-center space-x-1">
              {navItems.map((item) => {
                const Icon = item.icon
                return (
                  <Link key={item.path} to={item.path}>
                    <Button
                      variant={isActive(item.path) ? 'default' : 'ghost'}
                      className="flex items-center space-x-2"
                    >
                      <Icon className="h-4 w-4" />
                      <span>{item.label}</span>
                    </Button>
                  </Link>
                )
              })}
              {isHROrAdmin && (
                <Link to="/analytics">
                  <Button
                    variant={location.pathname.startsWith('/analytics') ? 'default' : 'ghost'}
                    className="flex items-center space-x-2"
                  >
                    <Settings className="h-4 w-4" />
                    <span>Analytics</span>
                  </Button>
                </Link>
              )}
              {isAuthenticated && (
                <div className="flex items-center space-x-2 ml-4 pl-4 border-l">
                  <span className="text-sm text-muted-foreground">{user?.username}</span>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={handleLogout}
                    className="flex items-center space-x-1"
                  >
                    <LogOut className="h-4 w-4" />
                    <span>Logout</span>
                  </Button>
                </div>
              )}
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <Outlet />
      </main>

      {/* Mobile Navigation */}
      <nav className="md:hidden fixed bottom-0 left-0 right-0 border-t bg-card p-2">
        <div className="flex justify-around">
          {navItems.slice(0, 3).map((item) => {
            const Icon = item.icon
            return (
              <Link key={item.path} to={item.path}>
                <Button
                  variant={isActive(item.path) ? 'default' : 'ghost'}
                  size="icon"
                  className="flex flex-col items-center space-y-1 h-auto py-2"
                >
                  <Icon className="h-5 w-5" />
                  <span className="text-xs">{item.label}</span>
                </Button>
              </Link>
            )
          })}
        </div>
      </nav>
    </div>
  )
}

export default Layout
