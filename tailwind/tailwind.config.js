/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/**/*.templ"],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        background: '#111111',
        primary: '#3498db',
        secondary: '#E74C3C',
        success: '#27AE60',
        warning: '#F39C12',
        muted: '#AAAAAA',
        headerFooterBg: '#222222',
        border: '#333333',
        codeHighlight: '#2ECC71',
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ]

}
