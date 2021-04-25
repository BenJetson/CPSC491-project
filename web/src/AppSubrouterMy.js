import React from "react";
import { Route, Switch, useRouteMatch } from "react-router-dom";
import NotFound from "./components/NotFound";
import NotImplemented from "./components/NotImplemented";

import { WithUser } from "./api/Auth";
import Roles from "./api/Roles";
import AccessDenied from "./components/AccessDenied";

import ProfileEditor from "./components/ProfileEditor";

const AppSubrouterMy = () => {
  const match = useRouteMatch();

  return (
    <Switch>
      <Route path={`${match.path}/profile`}>
        <ProfileEditor />
      </Route>
      <Route path={"*"}>
        {/* If no route matches, show a not found page. */}
        <NotFound />
      </Route>
    </Switch>
  );
};

export default AppSubrouterMy;
