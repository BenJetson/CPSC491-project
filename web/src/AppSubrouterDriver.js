import React from "react";
import { Route, Switch, useRouteMatch } from "react-router-dom";
import NotFound from "./components/NotFound";
import NotImplemented from "./components/NotImplemented";

import { WithUser } from "./api/Auth";
import Roles from "./api/Roles";
import AccessDenied from "./components/AccessDenied";
import ApplicationList from "./components/ApplicationList";
import ApplicationForm from "./components/ApplicationForm";

const AppSubrouterDriver = () => {
  const match = useRouteMatch();

  return (
    <WithUser>
      {({ isOneOfRoles }) =>
        (isOneOfRoles([Roles.IDOf.DRIVER]) && (
          <Switch>
            <Route exact path={`${match.path}/applications`}>
              <ApplicationList isSponsor={false} />
            </Route>
            <Route path={`${match.path}/applications/new`}>
              <ApplicationForm />
            </Route>
            <Route path={`${match.path}/applications/:id`}>
              <NotImplemented feature={"Driver - View Application"} />
            </Route>
            <Route path={`${match.path}/balance`}>
              <NotImplemented feature={"Driver - View Balance"} />
            </Route>
            <Route path={`${match.path}/shop`}>
              <NotImplemented feature={"Driver - Incentive Shop"} />
            </Route>
            <Route path={`${match.path}/receipts`}>
              <NotImplemented feature={"Driver - Receipts"} />
            </Route>
            <Route path={"*"}>
              {/* If no route matches, show a not found page. */}
              <NotFound />
            </Route>
          </Switch>
        )) || <AccessDenied />
      }
    </WithUser>
  );
};

export default AppSubrouterDriver;
