package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en" class="dark">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>APEX AI - Master the Future of AI Business</title>
	<!-- Google Fonts - Lexend Deca -->
	<link href="https://fonts.googleapis.com/css2?family=Lexend+Deca:wght@400;500;600;700&display=swap" rel="stylesheet">
	<!-- TailwindCSS (via CDN) -->
	<script src="https://cdn.tailwindcss.com"></script>
	<!-- Alpine.js -->
	<script defer src="https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js"></script>
	<!-- HTMX -->
	<script src="https://unpkg.com/htmx.org@1.9.2"></script>
	<!-- Matter.js -->
	<script src="https://cdnjs.cloudflare.com/ajax/libs/matter-js/0.19.0/matter.min.js"></script>
	<style>
		/* Apply Lexend Deca to all text */
		body {
			font-family: 'Lexend Deca', sans-serif;
		}
		/* Physics canvas container */
		#physics-container {
			position: fixed;
			width: 100%;
			height: 100vh;
			pointer-events: all;
			z-index: 1;
			transform-style: preserve-3d;
			will-change: transform;
			clip-path: none;
			top: 0;
			left: 0;
		}
		.logo-circle {
			width: 200px;
			height: 200px;
			border-radius: 50%;
			background: rgba(255, 255, 255, 0.15);
			backdrop-filter: blur(12px);
			display: flex;
			align-items: center;
			justify-content: center;
			font-size: 24px;
			color: white;
			text-align: center;
			line-height: 1.2;
			padding: 20px;
			box-shadow: 0 12px 24px rgba(0, 0, 0, 0.2);
		}
		/* Dynamic gradient animation */
		.gradient-animation {
			position: relative;
			overflow: hidden;
			background: #0a0a0a;
		}
		.gradient-layer {
			position: absolute;
			inset: -150%;
			opacity: 0.7;
			mix-blend-mode: plus-lighter;
			transition: transform 0.2s ease-out;
		}
		.gradient-layer-1 {
			background: radial-gradient(circle at center,
				transparent 0%,
				rgba(37, 99, 235, 0.9) 25%,  /* Bright blue */
				rgba(30, 58, 138, 0.8) 45%,   /* Dark blue */
				transparent 75%
			);
			animation: moveGradient1 25s linear infinite;
		}
		.gradient-layer-2 {
			background: radial-gradient(circle at center,
				transparent 0%,
				rgba(59, 130, 246, 0.8) 30%,   /* Medium blue */
				rgba(29, 78, 216, 0.7) 50%,   /* Blue 700 */
				transparent 80%
			);
			animation: moveGradient2 20s linear infinite;
		}
		.gradient-layer-3 {
			background: radial-gradient(circle at center,
				transparent 0%,
				rgba(96, 165, 250, 0.7) 35%,   /* Light blue */
				rgba(37, 99, 235, 0.6) 55%,     /* Blue 600 */
				transparent 85%
			);
			animation: moveGradient3 30s linear infinite;
		}
		/* Enhance animation ranges for more dynamic movement */
		@keyframes moveGradient1 {
			0% { transform: rotate(0deg) scale(1) translate(0%, 0%); }
			33% { transform: rotate(120deg) scale(1.4) translate(7%, 7%); }
			66% { transform: rotate(240deg) scale(0.8) translate(-7%, -7%); }
			100% { transform: rotate(360deg) scale(1) translate(0%, 0%); }
		}
		@keyframes moveGradient2 {
			0% { transform: rotate(360deg) scale(1.3) translate(7%, -7%); }
			33% { transform: rotate(240deg) scale(0.9) translate(-7%, 7%); }
			66% { transform: rotate(120deg) scale(1.2) translate(7%, 7%); }
			100% { transform: rotate(0deg) scale(1.3) translate(7%, -7%); }
		}
		@keyframes moveGradient3 {
			0% { transform: rotate(180deg) scale(0.8) translate(-7%, 7%); }
			33% { transform: rotate(300deg) scale(1.3) translate(7%, -7%); }
			66% { transform: rotate(60deg) scale(0.9) translate(-7%, -7%); }
			100% { transform: rotate(180deg) scale(0.8) translate(-7%, 7%); }
		}
		/* Adjust noise texture */
		.noise-overlay {
			position: absolute;
			inset: 0;
			background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%' height='100%' filter='url(%23noiseFilter)'/%3E%3C/svg%3E");
			opacity: 0.05;
			mix-blend-mode: color-dodge;
			pointer-events: none;
		}
		/* Make sure the solutions section has proper positioning */
		.solutions-section {
			position: relative;
			z-index: 1;
		}
	</style>
</head>
<body class="bg-gray-950 text-white antialiased">
	<!-- Physics container -->
	<div id="physics-container"></div>

	<!-- HERO SECTION -->
	<section class="gradient-animation min-h-[85vh] flex flex-col justify-center items-center text-center px-4 relative z-10"
		x-data="{ 
			mouseX: 0, 
			mouseY: 0,
			handleMouseMove(e) {
				const rect = e.currentTarget.getBoundingClientRect();
				this.mouseX = ((e.clientX - rect.left) / rect.width - 0.5) * 20;
				this.mouseY = ((e.clientY - rect.top) / rect.height - 0.5) * 20;
				this.$refs.layer1.style.transform += ' translate(' + (-this.mouseX) + 'px, ' + (-this.mouseY) + 'px)';
				this.$refs.layer2.style.transform += ' translate(' + (this.mouseX) + 'px, ' + (this.mouseY) + 'px)';
				this.$refs.layer3.style.transform += ' translate(' + (-this.mouseX * 0.5) + 'px, ' + (-this.mouseY * 0.5) + 'px)';
			}
		}"
		@mousemove="handleMouseMove"
	>
		<div class="gradient-layer gradient-layer-1" x-ref="layer1"></div>
		<div class="gradient-layer gradient-layer-2" x-ref="layer2"></div>
		<div class="gradient-layer gradient-layer-3" x-ref="layer3"></div>
		<div class="noise-overlay"></div>
		<h1 class="text-7xl md:text-9xl font-bold mb-4 tracking-tight relative z-10">
			<span class="italic">APEX</span> AI
		</h1>
		<div class="space-y-8 relative z-10">
			<!-- Main catchline -->
			<p class="text-4xl md:text-6xl lg:text-7xl font-medium bg-clip-text text-transparent bg-gradient-to-r from-gray-950 via-black to-gray-950">
				Turn AI into profit
			</p>
			<!-- Secondary catchlines -->
			<div class="max-w-2xl mx-auto space-y-1">
				<p class="text-lg md:text-xl text-gray-300 font-light tracking-wide">Your Proven Blueprint</p>
				<p class="text-base md:text-lg text-gray-400 font-light">for Business Growth</p>
			</div>
		</div>
		<!-- Scroll to form button -->
		<button onclick="document.getElementById('apply-form').scrollIntoView({behavior: 'smooth', block: 'center'})" 
			class="bg-blue-600 hover:bg-blue-700 transition px-8 py-4 rounded-lg text-lg relative z-10 mt-8">
			Apply Now
		</button>
	</section>

	<!-- BUBBLES SHOWCASE SECTION -->
	<section class="bubbles-section min-h-screen relative flex flex-col justify-center items-center py-32 px-4">
		<div class="max-w-4xl mx-auto text-center relative z-[2]">
			<h2 class="text-3xl md:text-4xl font-medium mb-8">Powered by Leading AI Technologies</h2>
			<p class="text-xl text-blue-200 mb-8">Leveraging the most advanced AI platforms to drive your business forward</p>
		</div>
		<div class="absolute inset-0 z-[1] bg-gradient-to-b from-gray-950/0 via-blue-950/10 to-gray-950/0"></div>
	</section>

	<!-- SOLUTION SECTION -->
	<section class="py-32 px-4 solutions-section min-h-screen relative z-10">
		<div class="max-w-3xl mx-auto">
			<h2 class="text-2xl font-medium mb-12 text-center">Solutions</h2>
			<div class="space-y-8">
				<div class="p-8 bg-blue-900/40 backdrop-blur-sm rounded-lg border border-blue-800/50">
					<h3 class="text-2xl font-medium mb-4">Expert Guidance</h3>
					<p class="text-blue-200 text-lg">Learn from industry experts with proven track records in AI implementation.</p>
				</div>
				<div class="p-8 bg-blue-900/40 backdrop-blur-sm rounded-lg border border-blue-800/50">
					<h3 class="text-2xl font-medium mb-4">Strategic Implementation</h3>
					<p class="text-blue-200 text-lg">Step-by-step blueprint for seamless AI integration into your business.</p>
				</div>
				<div class="p-8 bg-blue-900/40 backdrop-blur-sm rounded-lg border border-blue-800/50">
					<h3 class="text-2xl font-medium mb-4">ROI Optimization</h3>
					<p class="text-blue-200 text-lg">Proven frameworks to maximize your return on AI investments.</p>
				</div>
			</div>
		</div>
	</section>

	<!-- APPLICATION FORM SECTION -->
	<section id="apply-form" class="py-16 px-4">
		<div class="max-w-xl mx-auto">
			<h2 class="text-2xl font-medium mb-2 text-center">Apply Now</h2>
			<p class="text-gray-400 text-center mb-8">Join the ranks of successful AI-powered businesses</p>
			<div class="border border-blue-800/50 rounded-lg p-6">
				<form hx-post="/apply" hx-target="#form-response" hx-swap="innerHTML" class="space-y-4">
					<div>
						<label class="block mb-1 text-sm text-blue-200">Name</label>
						<input type="text" name="name" required 
							class="w-full p-2 rounded bg-blue-900/40 border border-blue-800/50 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition text-white">
					</div>
					<div>
						<label class="block mb-1 text-sm text-blue-200">Email</label>
						<input type="email" name="email" required 
							class="w-full p-2 rounded bg-blue-900/40 border border-blue-800/50 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition text-white">
					</div>
					<button type="submit" 
						class="w-full bg-blue-600 hover:bg-blue-700 transition px-4 py-2 rounded text-sm font-medium">
						Submit Application
					</button>
				</form>
				<div id="form-response" class="mt-4"></div>
			</div>
		</div>
	</section>

	<!-- PROOF SECTION -->
	<section class="py-16 px-4">
		<div class="max-w-2xl mx-auto">
			<h2 class="text-2xl font-medium mb-8 text-center">Testimonials</h2>
			<div x-data="{
				testimonials: [
					{ quote: 'This course transformed my business!', name: 'Alex' },
					{ quote: 'A must for anyone in AI.', name: 'Jordan' },
					{ quote: 'Unparalleled expertise and support.', name: 'Casey' }
				],
				current: 0
			}" class="border border-blue-800/50 rounded-lg p-6 text-center">
				<template x-if="testimonials.length">
					<div>
						<p class="text-lg text-blue-200" x-text="testimonials[current].quote"></p>
						<p class="mt-2 text-sm text-blue-300" x-text="testimonials[current].name"></p>
					</div>
				</template>
			</div>
		</div>
	</section>

	<!-- FINAL CTA SECTION -->
	<section class="py-16 px-4 border-t border-blue-800/50">
		<div class="max-w-xl mx-auto text-center">
			<h2 class="text-2xl font-medium mb-4">Ready to Start?</h2>
			<button onclick="document.getElementById('apply-form').scrollIntoView({behavior: 'smooth'})" 
				class="bg-blue-600 hover:bg-blue-700 transition px-6 py-2 rounded text-sm">
				Apply Now
			</button>
		</div>
	</section>

	<!-- Add this script before the closing body tag -->
	<script>
		// Initialize Matter.js
		const Engine = Matter.Engine,
			Render = Matter.Render,
			World = Matter.World,
			Bodies = Matter.Bodies,
			Body = Matter.Body,
			Mouse = Matter.Mouse,
			Events = Matter.Events;

		// Create engine and world
		const engine = Engine.create();
		const world = engine.world;
		
		// Adjust gravity and collision settings
		engine.world.gravity.y = 0.2;
		engine.timing.timeScale = 1;
		engine.enableSleeping = false;

		// Adjust collision detection settings for better precision
		const runner = Matter.Runner.create({
			isFixed: true,
			delta: 1000 / 60
		});

		// Create renderer with better settings
		const render = Render.create({
			element: document.getElementById('physics-container'),
			engine: engine,
			options: {
				width: window.innerWidth,
				height: window.innerHeight,
				wireframes: false,
				background: 'transparent',
				pixelRatio: window.devicePixelRatio
			}
		});

		// Company logos with their image URLs
		const logos = [
			{
				name: 'OpenAI',
				url: '/assets/logos/openai.png'
			},
			{
				name: 'Anthropic',
				url: '/assets/logos/anthropic.png'
			},
			{
				name: 'Grok',
				url: '/assets/logos/grok.png'
			},
			{
				name: 'Cohere',
				url: '/assets/logos/cohere.png'
			},
			{
				name: 'Stability',
				url: '/assets/logos/stability.png'
			},
			{
				name: 'DeepMind',
				url: '/assets/logos/deepmind.png'
			},
			{
				name: 'Perplexity',
				url: '/assets/logos/perplexity.png'
			},
			{
				name: 'Gemini',
				url: '/assets/logos/gemini.png'
			},
			{
				name: 'Mistral',
				url: '/assets/logos/mistral.png'
			}
		];

		// Function to create logo texture
		function createLogoTexture(img, name) {
			const canvas = document.createElement('canvas');
			const size = 180;
			canvas.width = size;
			canvas.height = size;
			const ctx = canvas.getContext('2d');

			// Draw transparent circle with white contour - this will be the exact physics boundary
			ctx.beginPath();
			ctx.arc(size/2, size/2, size/2 - 2, 0, Math.PI * 2);
			ctx.strokeStyle = 'rgba(255, 255, 255, 0.4)';
			ctx.lineWidth = 2;
			ctx.stroke();

			// Add subtle background for the circle
			ctx.beginPath();
			ctx.arc(size/2, size/2, size/2 - 3, 0, Math.PI * 2);
			ctx.fillStyle = 'rgba(255, 255, 255, 0.03)';
			ctx.fill();

			// Add glow effect
			ctx.shadowColor = 'rgba(255, 255, 255, 0.3)';
			ctx.shadowBlur = 15;
			ctx.shadowOffsetX = 0;
			ctx.shadowOffsetY = 0;

			// Create circular clip path matching the contour exactly
			ctx.save();
			ctx.beginPath();
			ctx.arc(size/2, size/2, size/2 - 4, 0, Math.PI * 2);
			ctx.clip();

			// Calculate dimensions to make logo smaller (60% of the circle size)
			const maxLogoSize = size * 0.6;
			const scale = Math.min(maxLogoSize / img.width, maxLogoSize / img.height);
			const width = img.width * scale;
			const height = img.height * scale;
			const x = size/2 - width/2;
			const y = size/2 - height/2;

			// Draw and process image
			ctx.drawImage(img, x, y, width, height);
			ctx.restore();

			// Add subtle inner glow
			ctx.beginPath();
			ctx.arc(size/2, size/2, size/2 - 3, 0, Math.PI * 2);
			ctx.strokeStyle = 'rgba(255, 255, 255, 0.1)';
			ctx.lineWidth = 4;
			ctx.stroke();

			return canvas.toDataURL();
		}

		// Create boundaries with adjusted properties
		const ground = Bodies.rectangle(
			window.innerWidth / 2,
			window.innerHeight + 60,
			window.innerWidth * 2,
			120,
			{ 
				isStatic: true, 
				render: { fillStyle: 'transparent' },
				restitution: 0.7,
				friction: 0.2
			}
		);

		const leftWall = Bodies.rectangle(
			-60,
			window.innerHeight / 2,
			120,
			window.innerHeight * 2,
			{ 
				isStatic: true, 
				render: { fillStyle: 'transparent' },
				restitution: 0.7,
				friction: 0.2
			}
		);

		const rightWall = Bodies.rectangle(
			window.innerWidth + 60,
			window.innerHeight / 2,
			120,
			window.innerHeight * 2,
			{ 
				isStatic: true, 
				render: { fillStyle: 'transparent' },
				restitution: 0.7,
				friction: 0.2
			}
		);

		// Add boundaries to world
		World.add(world, [ground, leftWall, rightWall]);

		// Shuffle and create bubbles
		const totalBubbles = logos.length;
		const columns = Math.ceil(Math.sqrt(totalBubbles)); // Calculate columns based on total bubbles
		const rows = Math.ceil(totalBubbles / columns);
		const columnWidth = window.innerWidth / columns;
		const rowHeight = 400;

		// Shuffle the logos array
		const shuffledLogos = [...logos];
		for (let i = shuffledLogos.length - 1; i > 0; i--) {
			const j = Math.floor(Math.random() * (i + 1));
			[shuffledLogos[i], shuffledLogos[j]] = [shuffledLogos[j], shuffledLogos[i]];
		}

		// Create bubbles using shuffled logos
		for (let i = 0; i < totalBubbles; i++) {
			const column = i % columns;
			const row = Math.floor(i / columns);
			const x = columnWidth * (column + 0.5) + (Math.random() - 0.5) * columnWidth * 0.3;
			const y = -rowHeight * row - Math.random() * 200;
			
			// Create bubble with specific logo instead of random
			const img = new Image();
			img.crossOrigin = "anonymous";
			img.src = shuffledLogos[i].url;
			
			img.onload = () => {
				const bubble = Bodies.circle(
					x,
					y,
					88,
					{
						restitution: 0.3,
						friction: 0.2,
						density: 0.01,
						frictionAir: 0.005,
						slop: 0,
						collisionFilter: {
							group: 0,
							category: 0x0001,
							mask: 0xFFFFFFFF
						},
						plugin: {
							attractors: [
								function(bodyA, bodyB) {
									const dx = bodyB.position.x - bodyA.position.x;
									const dy = bodyB.position.y - bodyA.position.y;
									const distance = Math.sqrt(dx * dx + dy * dy);
									
									if (distance < 88 * 3) {
										const force = -0.005 * Math.pow((1 - distance / (88 * 3)), 2);
										return {
											x: dx * force,
											y: dy * force
										};
									}
									return null;
								}
							]
						},
						render: {
							sprite: {
								texture: createLogoTexture(img, shuffledLogos[i].name),
								xScale: 1,
								yScale: 1
							}
						}
					}
				);
				
				Body.setVelocity(bubble, {
					x: (Math.random() - 0.5) * 1,
					y: 1
				});
				
				World.add(world, bubble);
			};
		}

		// Enhanced collision handling with stronger forces
		Events.on(engine, 'collisionStart', function(event) {
			event.pairs.forEach(function(pair) {
				if (!pair.bodyA.isStatic && !pair.bodyB.isStatic) {
					const dx = pair.bodyB.position.x - pair.bodyA.position.x;
					const dy = pair.bodyB.position.y - pair.bodyA.position.y;
					const distance = Math.sqrt(dx * dx + dy * dy);
					const force = 0.02; // Much stronger separation force
					
					Body.setVelocity(pair.bodyA, {
						x: pair.bodyA.velocity.x - (dx / distance) * force,
						y: pair.bodyA.velocity.y - (dy / distance) * force
					});
					
					Body.setVelocity(pair.bodyB, {
						x: pair.bodyB.velocity.x + (dx / distance) * force,
						y: pair.bodyB.velocity.y + (dy / distance) * force
					});

					// Add extra upward boost when colliding
					Body.setVelocity(pair.bodyA, {
						x: pair.bodyA.velocity.x,
						y: pair.bodyA.velocity.y - 1
					});
					Body.setVelocity(pair.bodyB, {
						x: pair.bodyB.velocity.x,
						y: pair.bodyB.velocity.y - 1
					});
				}
			});
		});

		// Continuous collision handling with stronger separation
		Events.on(engine, 'collisionActive', function(event) {
			event.pairs.forEach(function(pair) {
				if (!pair.bodyA.isStatic && !pair.bodyB.isStatic) {
					const dx = pair.bodyB.position.x - pair.bodyA.position.x;
					const dy = pair.bodyB.position.y - pair.bodyA.position.y;
					const distance = Math.sqrt(dx * dx + dy * dy);
					const minDistance = (pair.bodyA.circleRadius + pair.bodyB.circleRadius) * 1.5; // Increased minimum separation
					
					if (distance < minDistance) {
						const force = 0.01 * (minDistance - distance) / minDistance; // Stronger force
						Body.applyForce(pair.bodyA, pair.bodyA.position, {
							x: -dx * force,
							y: -dy * force
						});
						
						Body.applyForce(pair.bodyB, pair.bodyB.position, {
							x: dx * force,
							y: dy * force
						});

						// Add constant upward force when overlapping
						Body.applyForce(pair.bodyA, pair.bodyA.position, { x: 0, y: -0.0005 });
						Body.applyForce(pair.bodyB, pair.bodyB.position, { x: 0, y: -0.0005 });
					}
				}
			});
		});

		// Add click explosion effect with more controlled force
		render.canvas.addEventListener('click', (event) => {
			const rect = render.canvas.getBoundingClientRect();
			const clickX = event.clientX - rect.left;
			const clickY = event.clientY - rect.top;
			
			const bodies = Matter.Composite.allBodies(world);
			const explosionForce = 0.15;
			const explosionRadius = 150;

			bodies.forEach(body => {
				if (!body.isStatic) {
					const dx = body.position.x - clickX;
					const dy = body.position.y - clickY;
					const distance = Math.sqrt(dx * dx + dy * dy);

					if (distance < explosionRadius) {
						const force = (1 - distance / explosionRadius) * explosionForce;
						Body.applyForce(body, body.position, {
							x: (dx / distance) * force,
							y: (dy / distance) * force - force * 0.5 // Stronger upward boost
						});
						
						Body.setAngularVelocity(body, (Math.random() - 0.5) * 0.1); // Less rotation
					}
				}
			});
		});

		// Start the engine with the custom runner
		Matter.Runner.run(runner, engine);
		Render.run(render);

		// Handle window resize
		window.addEventListener('resize', () => {
			render.canvas.width = window.innerWidth;
			render.canvas.height = window.innerHeight;
		});
	</script>
</body>
</html>
`))

func main() {
	// Serve static files from the assets directory
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Serve the landing page.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	// Handle form submission via HTMX.
	http.HandleFunc("/apply", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// In a real-world application, process and store the form data here.
		name := r.FormValue("name")
		email := r.FormValue("email")
		log.Printf("Received application: %s <%s>", name, email)
		// Respond with a success snippet to be swapped into the page.
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="bg-green-500 p-4 rounded">Thank you for applying, ` + template.HTMLEscapeString(name) + `! We will be in touch shortly.</div>`))
	})

	log.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
