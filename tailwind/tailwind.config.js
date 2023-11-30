/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/**/*.templ", "./static/*.html"],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        background: '#111111',
        'steel-gray': {
          '50': '#1d1f20',
          '100': '#212325D',
          '200': '#282b2d',
          '300': '#323639',
          '400': '#1d2366',
          '500': '#242575',
          '600': '#342e83',
          '700': '#47408a',
          '800': '#3a3772',
          '900': '#32315b',
          '950': '#1d1b32',
        },
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ]

}
