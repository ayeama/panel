export const ENVIRONMENTS = [
  {
    name: 'evatt',
    host: 'panel.ayeama.com:8443/api',
  },
  // {
  //   name: 'eu-west-2',
  //   host: '127.0.0.1:8001',
  // },
]

export const HOST = ENVIRONMENTS[0].host

// TODO theme
// export function theme() {
//     var lightTheme = true

//     document.cookie = "theme=system"

//     if (true) {
//         lightTheme = window.matchMedia('(prefers-color-scheme: light)').matches
//         console.log(lightTheme)
//     }

//     if (lightTheme) {
//         document.documentElement.setAttribute('data-bs-theme', 'light')
//     } else {
//         document.documentElement.setAttribute('data-bs-theme', 'dark')
//     }
// }

// theme()

// window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
//     theme()
// })
