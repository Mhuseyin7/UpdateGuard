# UpdateGuard

> Safe update → verify → rollback if unhealthy.

UpdateGuard, Docker ve Docker Compose workload'ları için self-hosted bir **container update orchestration platform**'udur. Amaç, homelab kullanıcılarının, developer'ların ve küçük team'lerin image update işlemlerini yalnızca `docker pull` çalıştırarak değil; plan, health verification, stabilization ve gerektiğinde rollback ile güvenle yönetmesidir.

Bu açık kaynak kodlu software, [muhammedkoca.com.tr](https://muhammedkoca.com.tr) tarafından geliştirilmiştir.

## Neden UpdateGuard?

Bir image tag'inin değişmesi, deployment'ın başarılı olduğu anlamına gelmez. Yeni container açılmış görünebilir; ancak birkaç saniye sonra crash-loop'a girebilir, HTTP endpoint'i cevap vermeyebilir veya dependency bağlantısı kopabilir.

UpdateGuard'ın temel yaklaşımı şudur:

```text
Discover update
      ↓
Build exact plan
      ↓
Preflight + immutable snapshot
      ↓
Deploy candidate image
      ↓
Health verification + stabilization window
      ↓
SUCCESS  ← healthy
  or
ROLLBACK → previous known-good digest
```

Bu project bir Watchtower clone'u, otomatik `docker pull` cron job'u veya Docker socket'i browser'a açan bir dashboard değildir.

## Öne çıkan özellikler

- Safe update lifecycle ve doğrulanan state machine
- Mutable tag yerine digest-pinned snapshot ve rollback reference
- Docker socket'i yalnızca restricted agent'a ayıran architecture
- Docker HEALTHCHECK, HTTP/TCP/custom verification için model
- Stabilization window boyunca health, restart ve exit gözlemi
- Başarısız candidate için otomatik rollback workflow
- `ROLLBACK_FAILED` durumunu açıkça raporlama; yanlış başarı bilgisi yok
- Compose service scope'u dışında container recreate etmeme prensibi
- SemVer PATCH / MINOR / MAJOR classification (yalnızca gerçek SemVer tag'lerde)
- Infrastructure-focused Next.js dashboard shell
- Audit, policy, registry credential ve remote-agent security için tasarım dokümanları

## Update lifecycle

Update job'ları geçiş doğrulaması olan explicit state machine kullanır:

```text
PENDING
  → PRECHECK
  → BACKING_UP
  → PULLING
  → DEPLOYING
  → VERIFYING
  → STABILIZING
  → SUCCESS
```

Candidate unhealthy ise akış:

```text
VERIFYING / STABILIZING
  → ROLLING_BACK
  → ROLLED_BACK
```

Rollback da başarısız olursa final state `ROLLBACK_FAILED` olur. UpdateGuard recovery başarılıymış gibi davranmaz.

## Safety principles

- Volume'lar otomatik olarak silinmez.
- Rollback için gerekli previous image, update doğrulanana kadar prune edilmez.
- Mutable tag tek başına rollback source olarak kullanılmaz.
- Frontend veya public API Docker Engine'e direct access almaz.
- Registry password, API response, audit record veya log içinde saklanmaz.
- Major update'ler default olarak manual approval gerektirir.
- Unrestricted arbitrary shell command çalıştırılmaz.
- Aynı service/stack için concurrent update engellenmelidir.

## Architecture

```text
Browser
   │
   ▼
Web UI ──► Control Plane ──► Restricted Docker Agent ──► Docker Engine
                 │
                 └────────► Durable Store (SQLite / PostgreSQL)
```

Control plane policy, plan, job state ve audit ownership'ini taşır. Docker socket'i almaz. Docker Agent ise yalnızca discovery, image pull, service-scoped recreate, verification, rollback ve bounded logs için allowlisted operation'lara sahip olmalıdır.

Detaylar için [ARCHITECTURE.md](ARCHITECTURE.md) dosyasını inceleyin.

## Quick start

### Requirements

- Docker Engine ve Docker Compose v2
- Go 1.23+ (backend development için)
- Node.js 24+ ve npm (web development için)

### Local development

```bash
cp .env.example .env
# Set UPDATEGUARD_AGENT_TOKEN to a long, random value in .env.

go test ./...

cd web
npm ci
npm run typecheck
npm run build

cd ..
docker compose up --build -d
```

Ardından:

- Web UI: `http://localhost:3000`
- Control plane health: `http://localhost:8080/health`
- Control plane readiness: `http://localhost:8080/ready`

> Docker socket mount yalnızca agent service içindir. Agent port'unu internete publish etmeyin ve Docker TCP port'unu public olarak açmayın.

## Mevcut durum ve roadmap

Bu repository şu anda güvenlik odaklı bir foundation içerir: lifecycle, immutable plan rules, rollback orchestration, agent boundary, dashboard shell, testler ve deployment skeleton hazırdır.

Production workload yönetimi açılmadan önce aşağıdaki implementation milestone'ları tamamlanmalı ve review edilmelidir:

1. Docker Engine API ve Docker Compose discovery/recreate implementation
2. SQLite/PostgreSQL durable persistence, lock ve retention
3. Registry digest resolution ve private registry authentication
4. Authentication, RBAC, encrypted credentials ve audit persistence
5. mTLS remote-agent enrollment / rotation
6. Docker integration fixture: `v1 healthy → v2 unhealthy → rollback → v1 restored`

Bu yaklaşım bilinçlidir: eksik bir safety control varken container update'i "çalışıyor" olarak sunmak yerine, platform açıkça production-ready olmadığını belirtir.

## Verification

Frontend için:

```bash
cd web
npm run typecheck
npm run build
npm audit --omit=dev
```

Backend testleri:

```bash
go test ./...
go test -race ./...
```

CI pipeline Go formatting, vet, race tests, frontend typecheck/build ve Docker image build adımlarını içerir.

## Security

UpdateGuard high-trust bir infrastructure component'tir. Docker daemon erişimi host-level impact yaratabilir. Kurulumdan önce mutlaka aşağıdaki dokümanları okuyun:

- [SECURITY.md](SECURITY.md)
- [THREAT_MODEL.md](THREAT_MODEL.md)
- [ARCHITECTURE.md](ARCHITECTURE.md)

Security issue'ları public issue olarak paylaşmayın; [SECURITY.md](SECURITY.md) içindeki responsible disclosure yönergesini kullanın.

## Contributing

Contribution'lar memnuniyetle karşılanır. Her safety-sensitive değişiklik için test ekleyin; agent'a yeni bir Docker operation eklerken authorization test'i ve threat model güncellemesi zorunludur. Ayrıntılar: [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Bu repository için bir open-source license seçilmeden dağıtım koşulları tanımlı değildir. Yayına açmadan önce `LICENSE` dosyası ekleyin.

---

Built by [muhammedkoca.com.tr](https://muhammedkoca.com.tr) · Open-source container reliability tooling.
