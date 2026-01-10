/**
 * Code to manage the client-side session state, which is currently literally
 * just an auth token in localStorage. The flow is as follows: the token is set
 * on the login page; it is gotten in presumably every clientLoader in order to
 * make API requests; it is cleared in the API caller function golfFetch if it
 * ever produces a 401 Unauthorized response. TODO: add "log out" button
 * somewhere that also clears it
 */

export const golfAuthTokenKey = "golf:authToken";
export const setGolfAuthToken = (newToken: string) =>
  localStorage.setItem(golfAuthTokenKey, newToken);
export const getGolfAuthToken = () => localStorage.getItem(golfAuthTokenKey);
export const clearGolfAuthTOken = () => localStorage.removeItem(golfAuthTokenKey);
