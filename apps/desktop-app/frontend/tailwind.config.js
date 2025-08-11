/** @type {import('tailwindcss').Config} */
export default {
    darkMode: ["class"],
    content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}"

  ],
  theme: {
  	extend: {
  		fontFamily: {
  			title: [
  				'Sora',
  				'Arial',
  				'sans-serif'
  			],
  			text: [
  				'Poppins',
  				'Helvetica',
  				'Arial',
  				'sans-serif'
  			]
  		},
  		fontSize: {
  			'title-h1': '4rem',
  			'title-h2': '3.5rem',
  			'title-h3': '3rem',
  			'title-h4': '2.5rem',
  			'title-h5': '2.25rem',
  			'title-h6': '2rem',
  			'title-h7': '1.5rem',
  			'title-h8': '1.125rem',
  			'text-pl': '1.5rem',
  			'text-p2': '1.25rem',
  			'text-p3': '1.125rem',
  			'text-p4': '1rem',
  			'text-p5': '0.875rem',
  			'text-p6': '0.75rem',
  			'text-p7': '0.6875rem'
  		},
  		colors: {
			background: {
				'200': '#0B0E0F',
				'300': '#1A1B1F',
				'400': '#A786FD',
				'500': '#8A5EFC',
				'600': '#6D36FB',
				'700': '#572BC9',
				'800': '#412097',
				'900': '#2C1664',
				'1000': '#160B32'
			},
  			primary: {
  				'200': '#E2D7FE',
  				'300': '#C5AFFD',
  				'400': '#A786FD',
  				'500': '#8A5EFC',
  				'600': '#6D36FB',
  				'700': '#572BC9',
  				'800': '#412097',
  				'900': '#2C1664',
  				'1000': '#160B32'
  			},
  			secondary: {
  				'100': '#F4F0FF',
  				'200': '#EBFFD6',
  				'300': '#D6FFAD',
  				'400': '#C2FF85',
  				'500': '#ADFF5C',
  				'600': '#99FF33',
  				'700': '#7ACC29',
  				'800': '#5C991F',
  				'900': '#3D6614',
  				'1000': '#1F330A'
  			},
  			auxiliary: {
  				red: {
  					'100': '#FFEFEF',
  					'200': '#FFD7D7',
  					'300': '#FFBBB8',
  					'400': '#FF9797',
  					'500': '#FE5F5F',
  					'600': '#FF2E2E',
  					'700': '#CA1515',
  					'800': '#9E0B0B',
  					'900': '#750000',
  					'1000': '#4E0000'
  				},
  				orange: {
  					'100': '#FFF3E0',
  					'200': '#FFE2B7',
  					'300': '#FFD493',
  					'400': '#FFC46B',
  					'500': '#FFAF36',
  					'600': '#FF9900',
  					'700': '#D88100',
  					'800': '#BA7002',
  					'900': '#9C5E00',
  					'1000': '#774700'
  				},
  				yellow: {
  					'100': '#FFF9DE',
  					'200': '#FFF3B5',
  					'300': '#FFED8C',
  					'400': '#FFE871',
  					'500': '#FFE24D',
  					'600': '#FFD600',
  					'700': '#DEBB00',
  					'800': '#B99C03',
  					'900': '#A18803',
  					'1000': '#7D6A00'
  				},
  				blue: {
  					'100': '#E9F4FF',
  					'200': '#C2E2FF',
  					'300': '#ADD6FF',
  					'400': '#96CAFF',
  					'500': '#68B5FF',
  					'600': '#3399FF',
  					'700': '#1D88F3',
  					'800': '#1571CD',
  					'900': '#0A5FB4',
  					'1000': '#094179'
  				},
  				green: {
  					'100': '#E5FFE1',
  					'200': '#BFFFB4',
  					'300': '#A4F497',
  					'400': '#8BE37C',
  					'500': '#72D761',
  					'600': '#5EC34D',
  					'700': '#488639',
  					'800': '#399929',
  					'900': '#297A1B',
  					'1000': '#185810'
  				},
  				pink: {
  					'100': '#FFE5F9',
  					'200': '#FFBCEE',
  					'300': '#FFA0E7',
  					'400': '#FF8AE1',
  					'500': '#FF68D9',
  					'600': '#FF3FCE',
  					'700': '#DF32B3',
  					'800': '#C22199',
  					'900': '#A61381',
  					'1000': '#79005B'
  				}
  			},
  			gradient: {
  				linear: {
  					purple: '#F4F0FF',
  					green: '#E5FFE1',
  					red: '#FFEFEF',
  					pink: '#FFE5F9',
  					orange: '#FFF3E0',
  					yellow: '#FFF9DE',
  					blue: '#E9F4FF',
  					neutral: '#FFFFFF'
  				}
  			},
  			neutral: {
  				dark: {
  					'0': '#0B0B0E',
  					'100': '#0D0D12',
  					'200': '#19191E',
  					'300': '#252529',
  					'400': '#343438',
  					'500': '#414146',
  					'600': '#505057',
  					'700': '#585861',
  					'800': '#6D6D76',
  					'900': '#797983',
  					'1000': '#888898',
  					'1100': '#9D9DAB'
  				},
  				light: {
  					'0': '#FFFFFF',
  					'100': '#F7F7FC',
  					'200': '#EBEBF8',
  					'300': '#DCDC F0',
  					'400': '#BEBEDD',
  					'500': '#A9A9CC',
  					'600': '#949488',
  					'700': '#8383A4',
  					'800': '#777799',
  					'900': '#626284',
  					'1000': '#434360',
  					'1100': '#28283B'
  				}
  			},
  			sidebar: {
  				DEFAULT: 'hsl(var(--sidebar-background))',
  				foreground: 'hsl(var(--sidebar-foreground))',
  				primary: 'hsl(var(--sidebar-primary))',
  				'primary-foreground': 'hsl(var(--sidebar-primary-foreground))',
  				accent: 'hsl(var(--sidebar-accent))',
  				'accent-foreground': 'hsl(var(--sidebar-accent-foreground))',
  				border: 'hsl(var(--sidebar-border))',
  				ring: 'hsl(var(--sidebar-ring))'
  			}
  		},
  		borderRadius: {
  			lg: 'var(--radius)',
  			md: 'calc(var(--radius) - 2px)',
  			sm: 'calc(var(--radius) - 4px)'
  		}
  	}
  },
  plugins: [require("tailwindcss-animate")],
};
