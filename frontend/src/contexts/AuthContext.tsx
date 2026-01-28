import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { authApi, tokenStorage, userStorage, LoginResponse } from '@/services/auth';

interface AuthContextType {
  user: LoginResponse['user'] | null;
  isAuthenticated: boolean;
  isHROrAdmin: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  loading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<LoginResponse['user'] | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Load user from storage on mount
    const storedUser = userStorage.get();
    const token = tokenStorage.get();
    if (storedUser && token) {
      setUser(storedUser);
    }
    setLoading(false);
  }, []);

  const login = async (username: string, password: string) => {
    const response = await authApi.login({ username, password });
    tokenStorage.set(response.token);
    userStorage.set(response.user);
    setUser(response.user);
  };

  const logout = () => {
    tokenStorage.remove();
    userStorage.remove();
    setUser(null);
  };

  const isAuthenticated = !!user;
  const isHROrAdmin = user?.role === 'HR' || user?.role === 'ADMIN';

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated,
        isHROrAdmin,
        login,
        logout,
        loading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
