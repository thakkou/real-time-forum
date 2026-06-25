export const isAuthenticated = async () => {
  try {
    const response = await fetch(`${env.serverUri}/me`, {
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    });

    const resp = await response.json();

    if (!response.ok || resp.status_code === 401) {
      window.profile = null;

      return {
        authenticated: false,
      };
    }


    window.profile = {
      id: resp.data.id,
      nickname: resp.data.nickname,
    };
    
    return resp.data;
  } catch (err) {
    console.error(err);

    window.profile = null;

    return {
      authenticated: false,
    };
  }
};

// if (
//     route.auth &&
//     !localStorage.getItem('token')
// ) {
//     return this.navigate('/');
// }

// isAuthenticated (original)
//===========================

//  try {
//       const response = await fetch(`${env.serverUri}/me`, {
//           credentials: 'include',
//           headers: {
//               "Content-Type": "application/json",
//           },
//       });

//       const resp = await response.json();
//       // 401 unauthorized is returning in error in console (xhr)
//       if (resp.status_code === 401) return {
//           authenticated: false
//       };
//       return resp.data;
//   } catch (err){
//       console.log(err)
//       return {
//           authenticated: false
//       };