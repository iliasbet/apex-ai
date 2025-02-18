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
		/* Custom Scrollbar Styles */
		::-webkit-scrollbar {
			width: 8px;
			background: transparent;
		}
		::-webkit-scrollbar-track {
			background: rgba(0, 0, 0, 0.05);
			backdrop-filter: blur(8px);
		}
		::-webkit-scrollbar-thumb {
			background: rgba(255, 255, 255, 0.05);
			border: 1px solid rgba(255, 255, 255, 0.02);
			border-radius: 4px;
			backdrop-filter: blur(8px);
		}
		::-webkit-scrollbar-thumb:hover {
			background: rgba(255, 255, 255, 0.08);
		}
		/* Firefox */
		* {
			scrollbar-width: thin;
			scrollbar-color: rgba(255, 255, 255, 0.05) rgba(0, 0, 0, 0.05);
		}
		/* Physics canvas container */
		#physics-container {
			position: fixed;
			width: 100%;
			height: 100vh;
			pointer-events: none;
			z-index: 1;
			transform-style: preserve-3d;
			will-change: transform;
			clip-path: none;
			top: 0;
			left: 0;
		}
		#physics-container canvas {
			pointer-events: auto;
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
			background: rgba(0, 0, 0, 0.95);
		}
		.gradient-layer {
			position: absolute;
			inset: -150%;
			opacity: 0.8;
			mix-blend-mode: screen;
			transition: transform 0.2s ease-out;
		}
		.gradient-layer-1 {
			background: radial-gradient(circle at center,
				transparent 0%,
				rgba(255, 0, 128, 0.8) 25%,  /* Hot pink */
				rgba(128, 0, 255, 0.8) 45%,   /* Deep purple */
				transparent 75%
			);
			animation: moveGradient1 25s linear infinite;
		}
		.gradient-layer-2 {
			background: radial-gradient(circle at center,
				transparent 0%,
				rgba(0, 255, 200, 0.8) 30%,   /* Turquoise */
				rgba(0, 128, 255, 0.8) 50%,   /* Electric blue */
				transparent 80%
			);
			animation: moveGradient2 20s linear infinite;
		}
		.gradient-layer-3 {
			background: radial-gradient(circle at center,
				transparent 0%,
				rgba(255, 64, 0, 0.8) 35%,    /* Bright orange */
				rgba(255, 0, 64, 0.8) 55%,    /* Deep red */
				transparent 85%
			);
			animation: moveGradient3 30s linear infinite;
		}
		/* Enhance animation ranges for more dynamic movement */
		@keyframes moveGradient1 {
			0% { transform: rotate(0deg) scale(1.2) translate(0%, 0%); }
			33% { transform: rotate(120deg) scale(1.8) translate(15%, 15%); }
			66% { transform: rotate(240deg) scale(1) translate(-15%, -15%); }
			100% { transform: rotate(360deg) scale(1.2) translate(0%, 0%); }
		}
		@keyframes moveGradient2 {
			0% { transform: rotate(360deg) scale(1.5) translate(15%, -15%); }
			33% { transform: rotate(240deg) scale(1.1) translate(-15%, 15%); }
			66% { transform: rotate(120deg) scale(1.6) translate(15%, 15%); }
			100% { transform: rotate(0deg) scale(1.5) translate(15%, -15%); }
		}
		@keyframes moveGradient3 {
			0% { transform: rotate(180deg) scale(1) translate(-15%, 15%); }
			33% { transform: rotate(300deg) scale(1.7) translate(15%, -15%); }
			66% { transform: rotate(60deg) scale(1.1) translate(-15%, -15%); }
			100% { transform: rotate(180deg) scale(1) translate(-15%, 15%); }
		}
		/* Adjust noise texture */
		.noise-overlay {
			position: absolute;
			inset: 0;
			background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%' height='100%' filter='url(%23noiseFilter)'/%3E%3C/svg%3E");
			opacity: 0.08;
			mix-blend-mode: overlay;
			pointer-events: none;
		}
		/* Make sure the solutions section has proper positioning */
		.solutions-section {
			position: relative;
			z-index: 2;
			background: rgba(0, 0, 0, 0.7);
		}
		/* Update other sections */
		.bubbles-section {
			position: relative;
			z-index: 2;
			background: rgba(0, 0, 0, 0.7);
		}
		#apply-form {
			background: rgba(0, 0, 0, 0.7);
		}
		.bg-gray-950 {
			background-color: #000000;
		}
		
		/* Remove the gradient overlay in bubbles section */
		.bubbles-section .bg-gradient-to-b {
			display: none;
		}
		/* Add styles for translucent buttons */
		.btn-translucent {
			background: rgba(255, 255, 255, 0.15);
			backdrop-filter: blur(12px);
			border: 2px solid rgba(255, 255, 255, 0.4);
			transition: all 0.3s ease;
			text-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
			box-shadow: 0 0 20px rgba(255, 255, 255, 0.2),
						inset 0 0 15px rgba(255, 255, 255, 0.1);
			font-weight: 600;
			letter-spacing: 0.5px;
		}
		.btn-translucent:hover {
			background: rgba(255, 255, 255, 0.25);
			border-color: rgba(255, 255, 255, 0.6);
			transform: translateY(-2px);
			box-shadow: 0 0 30px rgba(255, 255, 255, 0.3),
						inset 0 0 20px rgba(255, 255, 255, 0.2);
			text-shadow: 0 0 15px rgba(255, 255, 255, 0.7);
		}
		/* Add glowing text effect */
		.glow-text {
			color: black;
			text-shadow: 0 0 10px rgba(255, 255, 255, 0.3),
						0 0 20px rgba(255, 255, 255, 0.2),
						0 0 30px rgba(255, 255, 255, 0.1);
			background: linear-gradient(to right,
				#000 0%,
				#222 50%,
				#000 100%
			);
			-webkit-background-clip: text;
			background-clip: text;
			animation: glow 2s ease-in-out infinite alternate;
		}

		@keyframes glow {
			from {
				text-shadow: 0 0 10px rgba(255, 255, 255, 0.3),
							0 0 20px rgba(255, 255, 255, 0.2),
							0 0 30px rgba(255, 255, 255, 0.1);
			}
			to {
				text-shadow: 0 0 15px rgba(255, 255, 255, 0.4),
							0 0 25px rgba(255, 255, 255, 0.3),
							0 0 35px rgba(255, 255, 255, 0.2);
			}
		}
	</style>
</head>
<body class="bg-gray-950 text-white antialiased">
	<!-- Physics container -->
	<div id="physics-container"></div>

	<!-- HERO SECTION -->
	<section class="gradient-animation min-h-[65vh] flex flex-col justify-center items-center text-center px-4 relative z-10"
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
		<h1 class="text-7xl md:text-9xl font-bold mb-8 tracking-tight relative z-10">
			<span class="italic text-blue-50/90">APEX</span> AI
		</h1>
		<div class="space-y-12 relative z-10">
			<!-- Main catchline -->
			<p class="text-3xl md:text-5xl lg:text-6xl font-medium text-black">
				Turn AI into profit
			</p>
		</div>
		<!-- Scroll to form button -->
		<button onclick="document.getElementById('apply-form').scrollIntoView({behavior: 'smooth', block: 'center'})" 
			class="btn-translucent px-10 py-5 rounded-lg text-xl relative z-10 mt-24 uppercase tracking-wider">
			Apply Now
		</button>
	</section>

	<!-- BUBBLES SHOWCASE SECTION -->
	<section class="bubbles-section min-h-screen relative flex flex-col justify-start items-center py-16 px-4">
		<div class="max-w-4xl mx-auto text-center relative z-10 mt-64">
			<h2 class="text-3xl md:text-4xl font-medium mb-4">Powered by Leading AI Technologies</h2>
			<p class="text-xl text-blue-200">Leveraging the most advanced AI platforms to drive your business forward</p>
		</div>
		<div class="absolute inset-0 z-[1] bg-gradient-to-b from-gray-950/0 via-blue-950/10 to-gray-950/0"></div>
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
			MouseConstraint = Matter.MouseConstraint,
			Events = Matter.Events;

		// Create engine and world
		const engine = Engine.create();
		const world = engine.world;
		
		// Adjust gravity and collision settings
		engine.world.gravity.y = 0.4;
		engine.timing.timeScale = 1.2;
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

		// Add mouse control
		const mouse = Mouse.create(render.canvas);
		const mouseConstraint = MouseConstraint.create(engine, {
			mouse: mouse,
			constraint: {
				stiffness: 0.5,
				render: {
					visible: false
				}
			},
			collisionFilter: {
				mask: 0x0001
			}
		});

		// Sync mouse with renderer
		render.mouse = mouse;

		// Add mouse constraint to world
		World.add(world, mouseConstraint);

		// Prevent page scrolling when dragging
		render.canvas.addEventListener('mousewheel', function(event) {
			if (mouseConstraint.body) {
				event.preventDefault();
			}
		});

		// Prevent default touch behavior when interacting with canvas
		render.canvas.addEventListener('touchmove', function(event) {
			if (mouseConstraint.body) {
				event.preventDefault();
			}
		}, { passive: false });

		// Add hover effect for bubbles
		Events.on(mouseConstraint, 'mousemove', function(event) {
			const mousePosition = event.mouse.position;
			const bodies = Matter.Composite.allBodies(world);
			
			bodies.forEach(body => {
				if (!body.isStatic) {
					const distance = Matter.Vector.magnitude(Matter.Vector.sub(body.position, mousePosition));
					if (distance < 100) {
						Body.setAngularVelocity(body, body.angularVelocity * 0.8);
						Body.setVelocity(body, {
							x: body.velocity.x * 0.8,
							y: body.velocity.y * 0.8
						});
					}
				}
			});
		});

		// Sync cursor style with dragging state
		Events.on(mouseConstraint, 'mousedown', function(event) {
			const mousePosition = event.mouse.position;
			const bodies = Matter.Composite.allBodies(world);
			let hovering = false;
			
			bodies.forEach(body => {
				if (!body.isStatic) {
					const distance = Matter.Vector.magnitude(Matter.Vector.sub(body.position, mousePosition));
					if (distance < 100) {
						hovering = true;
					}
				}
			});
			
			if (hovering) {
				render.canvas.style.cursor = 'grabbing';
			}
		});

		Events.on(mouseConstraint, 'mouseup', function(event) {
			render.canvas.style.cursor = 'default';
		});

		Events.on(mouseConstraint, 'mousemove', function(event) {
			const mousePosition = event.mouse.position;
			const bodies = Matter.Composite.allBodies(world);
			let hovering = false;
			
			bodies.forEach(body => {
				if (!body.isStatic) {
					const distance = Matter.Vector.magnitude(Matter.Vector.sub(body.position, mousePosition));
					if (distance < 100) {
						hovering = true;
					}
				}
			});
			
			render.canvas.style.cursor = hovering ? 'grab' : 'default';
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

			// Draw circle with white contour
			ctx.beginPath();
			ctx.arc(size/2, size/2, size/2 - 2, 0, Math.PI * 2);
			ctx.strokeStyle = 'rgba(255, 255, 255, 0.8)';
			ctx.lineWidth = 2;
			ctx.stroke();

			// Remove the solid background fill
			// Add subtle ring effect instead
			ctx.beginPath();
			ctx.arc(size/2, size/2, size/2 - 3, 0, Math.PI * 2);
			ctx.strokeStyle = 'rgba(255, 255, 255, 0.2)';
			ctx.lineWidth = 1;
			ctx.stroke();

			// Add subtle glow effect
			ctx.shadowColor = 'rgba(255, 255, 255, 0.5)';
			ctx.shadowBlur = 15;
			ctx.shadowOffsetX = 0;
			ctx.shadowOffsetY = 0;

			// Create circular clip path
			ctx.save();
			ctx.beginPath();
			ctx.arc(size/2, size/2, size/2 - 4, 0, Math.PI * 2);
			ctx.clip();

			// Calculate dimensions for logo
			const maxLogoSize = size * 0.6;
			const scale = Math.min(maxLogoSize / img.width, maxLogoSize / img.height);
			const width = img.width * scale;
			const height = img.height * scale;
			const x = size/2 - width/2;
			const y = size/2 - height/2;

			// Draw image
			ctx.drawImage(img, x, y, width, height);
			ctx.restore();

			// Add subtle inner glow
			ctx.beginPath();
			ctx.arc(size/2, size/2, size/2 - 3, 0, Math.PI * 2);
			ctx.strokeStyle = 'rgba(255, 255, 255, 0.3)';
			ctx.lineWidth = 2;
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
		const rowHeight = 100; // Reduced from 200 to 100

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
			const y = -rowHeight * row - Math.random() * 50; // Reduced from 100 to 50
			
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
					x: (Math.random() - 0.5) * 2,
					y: 2
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

	log.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
