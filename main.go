package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en" class="dark scroll-smooth">
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
	<!-- Add Cooper Hewitt font -->
	<link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:ital,wght@0,400;0,500;0,600;1,400;1,500;1,600&display=swap" rel="stylesheet">
	<style>
		/* Apply Lexend Deca to all text */
		body {
			font-family: 'Lexend Deca', sans-serif;
			scroll-behavior: smooth;
		}

		/* Hide scrollbar but keep functionality */
		* {
			-ms-overflow-style: none !important;  /* IE and Edge */
			scrollbar-width: none !important;  /* Firefox */
		}

		*::-webkit-scrollbar {
			display: none !important;  /* Chrome, Safari and Opera */
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

		#physics-world {
			position: absolute;
			width: 100%;
			height: 100%;
		}
		#physics-container canvas {
			pointer-events: none;
		}
		#physics-container canvas.interactive {
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
			background: url('/assets/background/background.png') center center/cover no-repeat;
		}
		/* Remove old gradient layers */
		.gradient-layer {
			display: none;
		}
		/* Adjust noise texture */
		.noise-overlay {
			position: absolute;
			inset: 0;
			mix-blend-mode: multiply;
			pointer-events: none;
		}
		/* Section styling */
		.solutions-section {
			position: relative;
			z-index: 2;
			background: rgba(0, 0, 0, 0.4);
			margin: 0; /* Remove any margin */
			border: none; /* Remove any borders */
		}
		/* Remove any gaps between sections */
		section + section {
			margin-top: 0;
			padding-top: 0;
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
			border: 2px solid rgba(255, 255, 255, 0.3);
			transition: all 0.3s ease;
			text-shadow: 0 0 5px rgba(255, 255, 255, 0.3);
			box-shadow: 0 0 10px rgba(255, 255, 255, 0.1),
						inset 0 0 8px rgba(255, 255, 255, 0.05);
			font-weight: 600;
			letter-spacing: 0.5px;
		}
		.btn-translucent:hover {
			background: rgba(255, 255, 255, 0.2);
			border-color: rgba(255, 255, 255, 0.4);
			transform: translateY(-2px);
			box-shadow: 0 0 15px rgba(255, 255, 255, 0.15),
						inset 0 0 10px rgba(255, 255, 255, 0.1);
			text-shadow: 0 0 8px rgba(255, 255, 255, 0.4);
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
		/* Add Cooper Hewitt font */
		@import url('https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:ital,wght@0,400;0,500;0,600;1,400;1,500;1,600&display=swap');
		
		.number-outline {
			font-family: 'IBM Plex Sans', sans-serif; /* Using IBM Plex Sans as it's similar to Cooper Hewitt */
			font-style: italic;
			font-weight: 500;
			-webkit-text-stroke: 4px #0066FF;
			paint-order: stroke fill;
		}
	</style>
</head>
<body class="bg-gray-950 text-white antialiased">
	<!-- Physics container -->
	<div id="physics-container"></div>

	<!-- HERO SECTION -->
	<section class="solutions-section gradient-animation min-h-[100vh] flex flex-col items-center text-center px-4 relative z-10">
		<div class="noise-overlay"></div>
		<div class="flex flex-col items-center mt-[15vh]">
			<h1 class="text-7xl md:text-9xl font-bold mb-8 tracking-tight relative z-10">
				<span class="italic" style="color: #0066FF">APEX</span> AI
			</h1>
			<div class="space-y-8 relative z-10">
				<p class="text-3xl md:text-5xl lg:text-6xl font-medium text-white">
					Turn AI into profit.
				</p>
			</div>
			<button onclick="document.getElementById('apply-form').scrollIntoView({behavior: 'smooth', block: 'center', duration: 1000})" 
				class="btn-translucent px-10 py-5 rounded-lg text-xl relative z-10 mt-16 uppercase tracking-wider hover:transform hover:translate-y-[-2px] transition-all duration-300">
				Apply Now
			</button>
			<div class="flex flex-col items-center gap-4 relative z-10 mt-16">
				<p class="text-white text-lg uppercase tracking-wider font-medium" style="text-shadow: 0 2px 4px rgba(0, 0, 0, 0.8), 0 4px 12px rgba(0, 0, 0, 0.9)">See More</p>
				<div class="animate-bounce pointer-events-none">
					<svg class="w-12 h-12 text-white" style="filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.8)) drop-shadow(0 4px 12px rgba(0, 0, 0, 0.9))" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3"></path>
					</svg>
				</div>
			</div>
		</div>
	</section>

	<!-- BUSINESS TRANSFORMATION SECTION -->
	<section class="solutions-section min-h-[60vh] relative flex flex-col justify-center items-center py-8 px-4">
		<div class="max-w-4xl mx-auto text-center relative z-10">
			<h2 class="text-3xl md:text-4xl font-medium mb-3">Transform Your Business with AI</h2>
			<p class="text-xl text-blue-200">Unlock unprecedented growth and efficiency with our end-to-end AI transformation solutions</p>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-16">
				<div class="p-5 rounded-lg backdrop-blur-md bg-white/5 relative mt-24">
					<div class="absolute -top-32 left-1/2 -translate-x-1/2 text-[140px] text-transparent number-outline">1</div>
					<h3 class="text-xl font-medium mb-2 mt-8">Strategic Planning</h3>
					<p class="text-blue-200/90">Create a tailored AI roadmap that aligns with your business objectives, market position, and growth targets. We analyze your data infrastructure, identify opportunities, and design a comprehensive implementation strategy.</p>
				</div>
				<div class="p-5 rounded-lg backdrop-blur-md bg-white/5 relative mt-24">
					<div class="absolute -top-32 left-1/2 -translate-x-1/2 text-[140px] text-transparent number-outline">2</div>
					<h3 class="text-xl font-medium mb-2 mt-8">Implementation</h3>
					<p class="text-blue-200/90">Seamlessly integrate cutting-edge AI solutions into your workflow with our expert team. We handle everything from model development and testing to deployment and team training, ensuring minimal disruption.</p>
				</div>
				<div class="p-5 rounded-lg backdrop-blur-md bg-white/5 relative mt-24">
					<div class="absolute -top-32 left-1/2 -translate-x-1/2 text-[140px] text-transparent number-outline">3</div>
					<h3 class="text-xl font-medium mb-2 mt-8">Optimization</h3>
					<p class="text-blue-200/90">Continuously refine and scale your AI capabilities through advanced analytics, performance monitoring, and iterative improvements. We ensure your AI systems evolve with your business needs and market demands.</p>
				</div>
			</div>
		</div>
	</section>

	<!-- ROI METRICS SECTION -->
	<section class="solutions-section min-h-[60vh] relative flex flex-col justify-center items-center py-8 px-4">
		<div class="max-w-4xl mx-auto text-center relative z-10">
			<h2 class="text-3xl md:text-4xl font-medium mb-3">Proven Business Impact</h2>
			<p class="text-xl text-blue-200">Our clients achieve exceptional results through strategic AI implementation</p>
			<div class="grid grid-cols-1 md:grid-cols-4 gap-6 mt-8">
				<div class="p-5 rounded-lg backdrop-blur-md bg-white/5">
					<h3 class="text-3xl font-medium mb-2" style="color: #0066FF">40%</h3>
					<p class="text-blue-200/90">Average Cost Reduction in Operations</p>
				</div>
				<div class="p-5 rounded-lg backdrop-blur-md bg-white/5">
					<h3 class="text-3xl font-medium mb-2" style="color: #0066FF">5x</h3>
					<p class="text-blue-200/90">Increase in Team Productivity</p>
				</div>
				<div class="p-5 rounded-lg backdrop-blur-md bg-white/5">
					<h3 class="text-3xl font-medium mb-2" style="color: #0066FF">60%</h3>
					<p class="text-blue-200/90">Faster Strategic Decision-Making</p>
				</div>
				<div class="p-5 rounded-lg backdrop-blur-md bg-white/5">
					<h3 class="text-3xl font-medium mb-2" style="color: #0066FF">24/7</h3>
					<p class="text-blue-200/90">Automated AI-Powered Operations</p>
				</div>
			</div>
		</div>
	</section>

	<!-- TRUST SECTION -->
	<section class="solutions-section min-h-[60vh] relative flex flex-col justify-center items-center py-8 px-4">
		<div class="max-w-4xl mx-auto text-center relative z-10">
			<h2 class="text-3xl md:text-4xl font-medium mb-3">Trusted by Global Industry Leaders</h2>
			<p class="text-xl text-blue-200">Join the ranks of forward-thinking companies transforming their industries with APEX AI</p>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-8">
				<div class="p-6 rounded-lg backdrop-blur-md bg-white/5 text-left">
					<div class="flex items-start mb-4">
						<div class="flex-1">
							<h3 class="text-lg font-medium text-white">Sarah Chen</h3>
							<p class="text-sm text-blue-200">Chief Technology Officer, TechVision Global</p>
							<div class="flex items-center mt-1">
								<span class="text-yellow-400">★★★★★</span>
							</div>
						</div>
						<div class="text-3xl text-blue-200/50">"</div>
					</div>
					<p class="text-blue-200/90 italic">APEX AI revolutionized our data analytics capabilities, delivering insights that were previously unattainable. Their expertise in implementing AI solutions has given us a significant competitive advantage in our market.</p>
				</div>
				<div class="p-6 rounded-lg backdrop-blur-md bg-white/5 text-left">
					<div class="flex items-start mb-4">
						<div class="flex-1">
							<h3 class="text-lg font-medium text-white">Marcus Rodriguez</h3>
							<p class="text-sm text-blue-200">Chief Executive Officer, InnovateCorp</p>
							<div class="flex items-center mt-1">
								<span class="text-yellow-400">★★★★★</span>
							</div>
						</div>
						<div class="text-3xl text-blue-200/50">"</div>
					</div>
					<p class="text-blue-200/90 italic">The ROI we've achieved since partnering with APEX AI has exceeded all expectations. Their strategic guidance and technical expertise have been instrumental in our digital transformation journey.</p>
				</div>
				<div class="p-6 rounded-lg backdrop-blur-md bg-white/5 text-left">
					<div class="flex items-start mb-4">
						<div class="flex-1">
							<h3 class="text-lg font-medium text-white">Emily Zhang</h3>
							<p class="text-sm text-blue-200">Head of Innovation, FutureScale Technologies</p>
							<div class="flex items-center mt-1">
								<span class="text-yellow-400">★★★★★</span>
							</div>
						</div>
						<div class="text-3xl text-blue-200/50">"</div>
					</div>
					<p class="text-blue-200/90 italic">Working with APEX AI has accelerated our AI adoption timeline by years. Their comprehensive implementation strategy and hands-on approach ensured a smooth transition to AI-powered operations.</p>
				</div>
				<div class="p-6 rounded-lg backdrop-blur-md bg-white/5 text-left">
					<div class="flex items-start mb-4">
						<div class="flex-1">
							<h3 class="text-lg font-medium text-white">David Park</h3>
							<p class="text-sm text-blue-200">Director of AI Operations, DataFlow Systems</p>
							<div class="flex items-center mt-1">
								<span class="text-yellow-400">★★★★★</span>
							</div>
						</div>
						<div class="text-3xl text-blue-200/50">"</div>
					</div>
					<p class="text-blue-200/90 italic">APEX AI doesn't just implement solutions; they become true partners in your AI journey. Their deep technical knowledge and strategic insights have transformed our business operations and market position.</p>
				</div>
			</div>
		</div>
	</section>

	<!-- EXPERTISE SECTION -->
	<section class="solutions-section min-h-[60vh] relative flex flex-col justify-center items-center py-8 px-4">
		<div class="max-w-4xl mx-auto text-center relative z-10">
			<h2 class="text-3xl md:text-4xl font-medium mb-3">World-Class Expertise</h2>
			<p class="text-xl text-blue-200">Backed by a team of industry veterans and AI specialists</p>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-8">
				<div class="p-6 rounded-lg backdrop-blur-md bg-white/5">
					<h3 class="text-2xl font-medium mb-3">Technical Excellence</h3>
					<ul class="text-left text-blue-200/90 space-y-2">
						<li>• State-of-the-Art AI Model Development</li>
						<li>• Enterprise-Grade Solution Architecture</li>
						<li>• Cloud-Native Infrastructure Design</li>
						<li>• Advanced Security & Compliance</li>
					</ul>
				</div>
				<div class="p-6 rounded-lg backdrop-blur-md bg-white/5">
					<h3 class="text-2xl font-medium mb-3">Business Acumen</h3>
					<ul class="text-left text-blue-200/90 space-y-2">
						<li>• Data-Driven ROI Optimization</li>
						<li>• Custom Industry Solutions</li>
						<li>• Enterprise Change Management</li>
						<li>• Long-Term Growth Strategy</li>
					</ul>
				</div>
			</div>
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
			MouseConstraint = Matter.MouseConstraint;

		// Company logos
		const logos = [
			{ name: 'OpenAI', url: '/assets/logos/openai.png' },
			{ name: 'Anthropic', url: '/assets/logos/anthropic.png' },
			{ name: 'DeepMind', url: '/assets/logos/deepmind.png' },
			{ name: 'Cohere', url: '/assets/logos/cohere.png' },
			{ name: 'Grok', url: '/assets/logos/grok.png' },
			{ name: 'Stability', url: '/assets/logos/stability.png' },
			{ name: 'Mistral', url: '/assets/logos/mistral.png' },
			{ name: 'Gemini', url: '/assets/logos/gemini.png' },
			{ name: 'Perplexity', url: '/assets/logos/perplexity.png' },
			{ name: 'DeepSeek', url: '/assets/logos/deepseek.png' }
		];

		// Function to create logo texture
		function createLogoTexture(img) {
			const canvas = document.createElement('canvas');
			const size = 180;
			canvas.width = size;
			canvas.height = size;
			const ctx = canvas.getContext('2d');

			// Draw circle with white contour
			ctx.beginPath();
			ctx.arc(size/2, size/2, size/2 - 2, 0, Math.PI * 2);
			ctx.strokeStyle = 'rgba(255, 255, 255, 0.95)';
			ctx.lineWidth = 2;
			ctx.stroke();

			// Add subtle ring effect
			ctx.beginPath();
			ctx.arc(size/2, size/2, size/2 - 3, 0, Math.PI * 2);
			ctx.strokeStyle = 'rgba(255, 255, 255, 0.4)';
			ctx.lineWidth = 1;
			ctx.stroke();

			// Add glow effect
			ctx.shadowColor = 'rgba(255, 255, 255, 0.8)';
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

			// Add inner glow
			ctx.beginPath();
			ctx.arc(size/2, size/2, size/2 - 3, 0, Math.PI * 2);
			ctx.strokeStyle = 'rgba(255, 255, 255, 0.5)';
			ctx.lineWidth = 2;
			ctx.stroke();

			return canvas.toDataURL();
		}

		// Create engine and world
		const engine = Engine.create();
		const world = engine.world;
		
		// Adjust physics settings
		engine.world.gravity.y = 0.5;
		engine.timing.timeScale = 1.2;
		engine.enableSleeping = false;

		// Add scroll-based physics adjustment
		let lastScrollY = window.scrollY;
		let scrollVelocity = 0;
		const physicsContainer = document.getElementById('physics-container');

		// Create renderer
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

		// Create a container for the physics world
		const physicsWorld = document.createElement('div');
		physicsWorld.id = 'physics-world';
		physicsContainer.appendChild(physicsWorld);
		physicsWorld.appendChild(render.canvas);

		window.addEventListener('scroll', (e) => {
			// Calculate scroll velocity
			const currentScrollY = window.scrollY;
			scrollVelocity = (currentScrollY - lastScrollY) * 0.1;
			lastScrollY = currentScrollY;

			// Adjust gravity based on scroll direction and speed
			const scrollFactor = Math.min(Math.abs(scrollVelocity) * 0.1, 2);
			const baseGravity = 0.5;

			if (scrollVelocity > 0) {
				// Scrolling down: decrease gravity significantly to create elevator effect
				engine.world.gravity.y = baseGravity * (1 - scrollFactor);
			} else if (scrollVelocity < 0) {
				// Scrolling up: increase gravity for catch-up effect
				engine.world.gravity.y = baseGravity * (1 + scrollFactor * 1.5);
			} else {
				// Not scrolling: return to base gravity
				engine.world.gravity.y = baseGravity;
			}

			// Apply velocity to all bodies when scrolling
			const bodies = Matter.Composite.allBodies(engine.world);
			bodies.forEach(body => {
				if (!body.isStatic) {
					Matter.Body.setVelocity(body, {
						x: body.velocity.x,
						y: body.velocity.y - scrollVelocity * 0.15  // Increased effect for more pronounced elevator feeling
					});
				}
			});
		});

		// Smooth scroll velocity decay
		function updateScrollVelocity() {
			scrollVelocity *= 0.95;
			requestAnimationFrame(updateScrollVelocity);
		}
		updateScrollVelocity();

		// Create boundaries
		let boundaries = {
			ground: Bodies.rectangle(
				window.innerWidth / 2,
				window.innerHeight + 30,
				window.innerWidth,
				60,
				{ 
					isStatic: true,
					restitution: 0.9,    // High bounce for ground
					friction: 0.2,        // More friction for rolling
					render: { fillStyle: 'transparent' }
				}
			),
			ceiling: Bodies.rectangle(
				window.innerWidth / 2,
				-30,
				window.innerWidth,
				60,
				{ 
					isStatic: true,
					restitution: 0.9,    // High bounce for ceiling
					friction: 0.2,        // More friction for rolling
					render: { fillStyle: 'transparent' }
				}
			),
			leftWall: Bodies.rectangle(
				-30,
				window.innerHeight / 2,
				60,
				window.innerHeight,
				{ 
					isStatic: true,
					restitution: 0.9,    // High bounce for walls
					friction: 0.2,        // More friction for rolling
					render: { fillStyle: 'transparent' }
				}
			),
			rightWall: Bodies.rectangle(
				window.innerWidth + 30,
				window.innerHeight / 2,
				60,
				window.innerHeight,
				{ 
					isStatic: true,
					restitution: 0.9,    // High bounce for walls
					friction: 0.2,        // More friction for rolling
					render: { fillStyle: 'transparent' }
				}
			)
		};

		// Add boundaries to world
		World.add(world, Object.values(boundaries));

		// Calculate responsive bubble size based on viewport
		function calculateBubbleSize() {
			const minSize = 66; // Minimum bubble radius
			const maxSize = 88; // Maximum bubble radius
			const breakpoint = 1200; // Width at which we want max size
			
			const size = Math.min(maxSize, Math.max(minSize, 
				(window.innerWidth / breakpoint) * maxSize
			));
			
			return size;
		}

		// Create ball with realistic physics
		function createBubble(x, y, img) {
			const bubbleSize = calculateBubbleSize();
			return Bodies.circle(x, y, bubbleSize, {  
				restitution: 0.85,    // Realistic ball bounce
				friction: 0.1,        // Moderate friction for rolling
				frictionAir: 0.001,   // Some air resistance
				density: 0.002,       // More weight for momentum
				torque: 0.002,        // Add some spin
				render: {
					sprite: {
						texture: createLogoTexture(img),
						xScale: bubbleSize/88,     // Scale texture relative to original size
						yScale: bubbleSize/88      // Scale texture relative to original size
					}
				}
			});
		}

		// Create and add bubbles in a grid
		const totalBubbles = logos.length;
		const columns = Math.ceil(Math.sqrt(totalBubbles));
		const columnWidth = window.innerWidth / columns;

		// Load and create bubbles
		logos.forEach((logo, i) => {
			const column = i % columns;
			const row = Math.floor(i / columns);
			const x = columnWidth * (column + 0.5);
			const y = window.innerHeight * 0.4 + (100 * row); // Start bubbles at 40% of viewport height
			
			const img = new Image();
			img.crossOrigin = "anonymous";
			img.src = logo.url;
			
			img.onload = () => {
				const bubble = createBubble(x, y, img);
				World.add(world, bubble);

				// Initial velocity with more upward tendency
				Body.setVelocity(bubble, {
					x: (Math.random() - 0.5) * 4, // Gentle horizontal velocity
					y: Math.random() * -2 - 1  // Upward velocity
				});

				// Add some initial angular velocity for spin
				Body.setAngularVelocity(bubble, (Math.random() - 0.5) * 0.02); // Very gentle spin
			};
		});

		// Add mouse interaction
		const mouse = Mouse.create(render.canvas);
		const mouseConstraint = MouseConstraint.create(engine, {
			mouse: mouse,
			constraint: {
				stiffness: 0.2,
				render: { visible: false }
			}
		});

		World.add(world, mouseConstraint);
		render.mouse = mouse;

		// Make canvas interactive when hovering bubbles
		render.canvas.addEventListener('mousemove', () => {
			render.canvas.style.cursor = mouseConstraint.body ? 'grab' : 'default';
		});

		// Window resize handler
		window.addEventListener('resize', () => {
			render.canvas.width = window.innerWidth;
			render.canvas.height = window.innerHeight;
			
			// Remove old boundaries
			World.remove(world, Object.values(boundaries));
			
			// Update boundaries
			boundaries = {
				ground: Bodies.rectangle(
					window.innerWidth / 2,
					window.innerHeight + 30,
					window.innerWidth,
					60,
					{ 
						isStatic: true,
						restitution: 0.9,
						friction: 0.2,
						render: { fillStyle: 'transparent' }
					}
				),
				ceiling: Bodies.rectangle(
					window.innerWidth / 2,
					-30,
					window.innerWidth,
					60,
					{ 
						isStatic: true,
						restitution: 0.9,
						friction: 0.2,
						render: { fillStyle: 'transparent' }
					}
				),
				leftWall: Bodies.rectangle(
					-30,
					window.innerHeight / 2,
					60,
					window.innerHeight,
					{ 
						isStatic: true,
						restitution: 0.9,
						friction: 0.2,
						render: { fillStyle: 'transparent' }
					}
				),
				rightWall: Bodies.rectangle(
					window.innerWidth + 30,
					window.innerHeight / 2,
					60,
					window.innerHeight,
					{ 
						isStatic: true,
						restitution: 0.9,
						friction: 0.2,
						render: { fillStyle: 'transparent' }
					}
				)
			};
			
			// Add new boundaries
			World.add(world, Object.values(boundaries));

			// Update all bubbles with new size
			const bubbleSize = calculateBubbleSize();
			const bodies = Matter.Composite.allBodies(engine.world);
			bodies.forEach(body => {
				if (!body.isStatic) {
					const scale = bubbleSize / body.circleRadius;
					Matter.Body.scale(body, scale, scale);
					if (body.render.sprite) {
						body.render.sprite.xScale = bubbleSize/88;
						body.render.sprite.yScale = bubbleSize/88;
					}
				}
			});
		});

		// Start the engine
		Matter.Runner.run(engine);
		Render.run(render);
	</script>

	<!-- LEGAL FOOTER SECTION -->
	<footer class="solutions-section relative py-16 px-4 border-t border-white/10">
		<div class="absolute inset-x-0 bottom-0 h-64 bg-gradient-to-t from-black to-transparent"></div>
		<div class="max-w-7xl mx-auto p-5 rounded-lg backdrop-blur-md bg-white/5 relative z-[2]">
			<div class="grid grid-cols-1 md:grid-cols-4 gap-12 mb-12">
				<!-- Company Info -->
				<div>
					<h3 class="text-xl font-medium mb-4">APEX AI</h3>
					<p class="text-blue-200/90 text-sm">
						Leading the future of AI-driven business transformation.
					</p>
				</div>
				
				<!-- Quick Links -->
				<div>
					<h3 class="text-lg font-medium mb-4">Quick Links</h3>
					<ul class="space-y-2 text-sm">
						<li><a href="#" class="text-blue-200/90 hover:text-white transition">About Us</a></li>
						<li><a href="#" class="text-blue-200/90 hover:text-white transition">Services</a></li>
						<li><a href="#" class="text-blue-200/90 hover:text-white transition">Case Studies</a></li>
						<li><a href="#" class="text-blue-200/90 hover:text-white transition">Contact</a></li>
					</ul>
				</div>
				
				<!-- Legal -->
				<div>
					<h3 class="text-lg font-medium mb-4">Legal</h3>
					<ul class="space-y-2 text-sm">
						<li><a href="#" class="text-blue-200/90 hover:text-white transition">Terms of Service</a></li>
						<li><a href="#" class="text-blue-200/90 hover:text-white transition">Privacy Policy</a></li>
						<li><a href="#" class="text-blue-200/90 hover:text-white transition">Cookie Policy</a></li>
						<li><a href="#" class="text-blue-200/90 hover:text-white transition">GDPR Compliance</a></li>
					</ul>
				</div>
				
				<!-- Contact -->
				<div>
					<h3 class="text-lg font-medium mb-4">Contact Us</h3>
					<ul class="space-y-2 text-sm">
						<li class="text-blue-200/90">Email: contact@apexai.com</li>
						<li class="text-blue-200/90">Phone: +1 (555) 123-4567</li>
						<li class="text-blue-200/90">Location: San Francisco, CA</li>
					</ul>
				</div>
			</div>
			
			<!-- Bottom Bar -->
			<div class="pt-8 border-t border-white/10">
				<div class="flex flex-col md:flex-row justify-between items-center gap-4">
					<div class="text-sm text-blue-200/90">
						© 2024 APEX AI. All rights reserved.
					</div>
					<div class="flex gap-6">
						<a href="#" class="text-blue-200/90 hover:text-white transition">
							<span class="sr-only">LinkedIn</span>
							<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
								<path d="M19 0h-14c-2.761 0-5 2.239-5 5v14c0 2.761 2.239 5 5 5h14c2.762 0 5-2.239 5-5v-14c0-2.761-2.238-5-5-5zm-11 19h-3v-11h3v11zm-1.5-12.268c-.966 0-1.75-.79-1.75-1.764s.784-1.764 1.75-1.764 1.75.79 1.75 1.764-.783 1.764-1.75 1.764zm13.5 12.268h-3v-5.604c0-3.368-4-3.113-4 0v5.604h-3v-11h3v1.765c1.396-2.586 7-2.777 7 2.476v6.759z"/>
							</svg>
						</a>
						<a href="#" class="text-blue-200/90 hover:text-white transition">
							<span class="sr-only">Twitter</span>
							<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
								<path d="M23.953 4.57a10 10 0 01-2.825.775 4.958 4.958 0 002.163-2.723c-.951.555-2.005.959-3.127 1.184a4.92 4.92 0 00-8.384 4.482C7.69 8.095 4.067 6.13 1.64 3.162a4.822 4.822 0 00-.666 2.475c0 1.71.87 3.213 2.188 4.096a4.904 4.904 0 01-2.228-.616v.06a4.923 4.923 0 003.946 4.827 4.996 4.996 0 01-2.212.085 4.936 4.936 0 004.604 3.417 9.867 9.867 0 01-6.102 2.105c-.39 0-.779-.023-1.17-.067a13.995 13.995 0 007.557 2.209c9.053 0 13.998-7.496 13.998-13.985 0-.21 0-.42-.015-.63A9.935 9.935 0 0024 4.59z"/>
							</svg>
						</a>
						<a href="#" class="text-blue-200/90 hover:text-white transition">
							<span class="sr-only">GitHub</span>
							<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
								<path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd"/>
							</svg>
						</a>
					</div>
				</div>
				
				<!-- Legal Disclaimer -->
				<p class="mt-8 text-xs text-blue-200/70 text-center">
					Disclaimer: The information provided on this website is for general informational purposes only. While we strive to keep the information up to date and accurate, we make no representations or warranties of any kind, express or implied, about the completeness, accuracy, reliability, suitability, or availability with respect to the website or the information, products, services, or related graphics contained on the website for any purpose.
				</p>
			</div>
		</div>
	</footer>
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

	log.Println("Server started at http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
