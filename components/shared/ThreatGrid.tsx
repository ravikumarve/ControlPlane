"use client";

import { useEffect, useRef } from "react";

interface Node {
  x: number;
  y: number;
  type: "agent" | "tool";
  status: "nominal" | "alert";
  alertTimer: number;
  connections: Node[];
}

interface Packet {
  x: number;
  y: number;
  target: Node;
  progress: number;
  speed: number;
  isMalicious: boolean;
}

export function ThreatGrid() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    let width: number;
    let height: number;
    const nodes: Node[] = [];
    const packets: Packet[] = [];
    const numNodes = window.innerWidth > 768 ? 60 : 30;
    let animationId: number;

    function init() {
      width = canvas!.width = window.innerWidth;
      height = canvas!.height = window.innerHeight;

      nodes.length = 0;
      packets.length = 0;

      for (let i = 0; i < numNodes; i++) {
        nodes.push({
          x: Math.random() * width,
          y: Math.random() * height,
          type: Math.random() > 0.7 ? "tool" : "agent",
          status: "nominal",
          alertTimer: 0,
          connections: [],
        });
      }

      nodes.forEach((n) => {
        n.connections = [];
        nodes.forEach((n2) => {
          if (n !== n2 && Math.hypot(n.x - n2.x, n.y - n2.y) < 250) {
            if (
              (n.type === "agent" && n2.type === "tool") ||
              (n.type === "tool" && n2.type === "agent")
            ) {
              n.connections.push(n2);
            }
          }
        });
      });
    }

    function spawnPacket() {
      const agents = nodes.filter(
        (n) => n.type === "agent" && n.connections.length > 0
      );
      if (agents.length === 0) return;

      const start = agents[Math.floor(Math.random() * agents.length)];
      const end =
        start.connections[
          Math.floor(Math.random() * start.connections.length)
        ];

      packets.push({
        x: start.x,
        y: start.y,
        target: end,
        progress: 0,
        speed: 0.01 + Math.random() * 0.01,
        isMalicious: Math.random() > 0.95,
      });
    }

    function animate() {
      animationId = requestAnimationFrame(animate);
      ctx!.clearRect(0, 0, width, height);

      if (Math.random() > 0.8) spawnPacket();

      // Connections
      ctx!.lineWidth = 1;
      nodes.forEach((n) => {
        n.connections.forEach((n2) => {
          ctx!.beginPath();
          ctx!.moveTo(n.x, n.y);
          ctx!.lineTo(n2.x, n2.y);
          ctx!.strokeStyle = "rgba(255, 255, 255, 0.05)";
          ctx!.stroke();
        });
      });

      // Nodes
      nodes.forEach((n) => {
        ctx!.beginPath();
        if (n.alertTimer > 0) {
          n.alertTimer--;
          if (n.alertTimer <= 0) n.status = "nominal";
        }

        if (n.status === "alert") {
          ctx!.arc(n.x, n.y, 6, 0, Math.PI * 2);
          ctx!.fillStyle = "#f97316";
          ctx!.shadowBlur = 15;
          ctx!.shadowColor = "#f97316";
        } else if (n.type === "tool") {
          ctx!.rect(n.x - 3, n.y - 3, 6, 6);
          ctx!.fillStyle = "#06b6d4";
          ctx!.shadowBlur = 5;
          ctx!.shadowColor = "#06b6d4";
        } else {
          ctx!.arc(n.x, n.y, 2, 0, Math.PI * 2);
          ctx!.fillStyle = "rgba(255, 255, 255, 0.3)";
          ctx!.shadowBlur = 0;
        }

        ctx!.fill();
        ctx!.shadowBlur = 0;
      });

      // Packets
      for (let i = packets.length - 1; i >= 0; i--) {
        const p = packets[i];
        p.progress += p.speed;

        if (p.progress >= 1) {
          if (p.isMalicious) {
            p.target.status = "alert";
            p.target.alertTimer = 60;
          }
          packets.splice(i, 1);
          continue;
        }

        const cx = p.x + (p.target.x - p.x) * p.progress;
        const cy = p.y + (p.target.y - p.y) * p.progress;

        ctx!.beginPath();
        ctx!.arc(cx, cy, 2, 0, Math.PI * 2);
        ctx!.fillStyle = p.isMalicious ? "#f97316" : "#06b6d4";
        ctx!.shadowBlur = p.isMalicious ? 10 : 5;
        ctx!.shadowColor = p.isMalicious ? "#f97316" : "#06b6d4";
        ctx!.fill();
        ctx!.shadowBlur = 0;
      }
    }

    init();
    animate();

    const handleResize = () => init();
    window.addEventListener("resize", handleResize);

    return () => {
      cancelAnimationFrame(animationId);
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  return (
    <canvas
      ref={canvasRef}
      className="fixed left-0 top-0 h-screen w-screen opacity-60"
      style={{ zIndex: 0, pointerEvents: "none" }}
    />
  );
}
