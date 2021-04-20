import React, { useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { Request } from "./Base";

const AuthContext = React.createContext({ user: null });

const GetCurrentUser = async () => {
  const res = await Request("GET", "/whoami");

  if (res.error !== null) {
    return null;
  }
  return res.data;
};

const DoLogin = async (email, password) => {
  return await Request("POST", "/login", { email, password });
};

const DoLogout = async () => {
  return await Request("POST", "/logout");
};

const AuthProvider = ({ children }) => {
  const history = useHistory();

  const [user, setUser] = useState(null);
  const [lastCheck, setLastCheck] = useState(Date.now());

  useEffect(() => {
    (async () => {
      const currentUser = await GetCurrentUser();
      setUser(currentUser);
    })();
  }, [lastCheck]);

  const logout = () => {
    setUser(null);
    DoLogout();

    // Return to the login page.
    history.push("/login");
  };

  const refreshLoginStatus = () => {
    setLastCheck(Date.now());
  };

  const isAuthenticated = () => {
    return user !== null;
  };

  const isOneOfRoles = (roles) => {
    return roles.includes(user?.["role_id"] ?? false);
  };

  const getName = () => {
    const firstName = user?.["first_name"] ?? "Unknown";
    const lastName = user?.["last_name"] ?? "Unknown";

    return [firstName, lastName];
  };

  const getInitials = () => {
    const [firstName, lastName] = getName();
    const initials = (
      firstName.substring(0, 1) + lastName.substring(0, 1)
    ).toUpperCase();

    return initials;
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        logout,
        isAuthenticated,
        isOneOfRoles,
        refreshLoginStatus,
        getName,
        getInitials,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

const WithUser = AuthContext.Consumer;

export { AuthProvider, WithUser, DoLogin, DoLogout, GetCurrentUser };
