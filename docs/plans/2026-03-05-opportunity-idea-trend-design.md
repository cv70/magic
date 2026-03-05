# Opportunity Insight + Idea Factory + Trend Content Design

Date: 2026-03-05
Status: Approved
Scope: MVP

## 1. Context And Goals

This design extends the current AI content engine for two parallel outcomes on a shared foundation:
- Business opportunity insight
- Trend-aligned content generation with high growth potential

User-confirmed constraints:
- Primary target user for V1: individual creator
- Signal strategy: hybrid (platform trends + search demand + historical self data)
- Optimization objective: multi-objective with default weights
  - Viral potential: 40%
  - Follower growth: 30%
  - Business conversion: 30%
- Initial platforms: Xiaohongshu + WeChat Official Account

## 2. Product Architecture (MVP)

Use one insight foundation with two product lines:

1. Opportunity Insight Center
- Input: Xiaohongshu trend signals, WeChat keyword/search signals, historical account performance
- Output: opportunity cards (niche, audience, pain point, business potential, competition)

2. Idea Factory
- Input: opportunity cards
- Output: structured idea matrix (`topic x audience x angle x format`) with priority

3. Trend Content Generator
- Input: selected idea + target platform
- Output: platform-tailored draft (title, hook, structure, CTA) + 3-axis scores

4. Feedback Learning Loop
- Input: post-publish performance
- Output: recommendation tuning for subsequent opportunity ranking and content generation

System mapping to current codebase:
- Backend new domains: `insight`, `idea`, `trend_content`, `feedback`
- Frontend pages: insights board, idea workspace, content studio
- MVP scoring engine: rules + lightweight statistics

## 3. Core Data Model

### 3.1 Signal
- `source` (`xiaohongshu` | `wechat` | `search` | `self`)
- `topic`
- `keywords[]`
- `engagement_hint`
- `captured_at`

### 3.2 Opportunity
- `title`
- `audience`
- `pain_point`
- `commercial_intent`
- `competition_level`
- `novelty`
- `score_total`

### 3.3 Idea
- `opportunity_id`
- `angle`
- `format` (`list` | `tutorial` | `opinion` | `case`)
- `platform_fit`
- `estimated_scores`

### 3.4 GeneratedContent
- `idea_id`
- `platform`
- `title`
- `hook`
- `outline`
- `body`
- `cta`
- `score_breakdown`

### 3.5 PerformanceFeedback
- `content_id`
- `views`
- `reads`
- `likes`
- `saves`
- `follows`
- `leads`
- `published_at`

## 4. Scoring Framework

MVP uses deterministic scoring (rules first, model later):

- `viral_score`: trend intensity, growth velocity, hook alignment
- `growth_score`: persona consistency, follow intent quality, serial potential
- `conversion_score`: pain intensity, CTA clarity, lead intent strength

Global objective:

`total = 0.4 * viral + 0.3 * growth + 0.3 * conversion`

Platform-aware interpretation:
- Xiaohongshu: hook resonance + save/share tendency
- WeChat OA: structural completeness + trust depth + conversion path clarity

## 5. MVP Workflow

1. Signal ingestion + normalization
- Scheduled jobs ingest trend/search/self signals
- De-duplication and clustering into normalized signal records

2. Opportunity generation
- Build opportunity cards from normalized signal clusters
- Rank by comprehensive score and confidence

3. Idea expansion
- Expand each opportunity into multiple angles and formats
- Provide keep/drop/edit actions in idea workspace

4. Trend content generation
- Generate platform-specific drafts from selected ideas
- Return score breakdown + optimization hints

5. Post-publish feedback
- Manual feedback input first, API integration later
- Adjust ranking/generation weights via simple tuning logic

## 6. API Draft (MVP)

- `POST /api/v1/insights/signals/ingest`
- `GET /api/v1/insights/opportunities`
- `POST /api/v1/ideas/generate`
- `POST /api/v1/contents/generate-trend`
- `POST /api/v1/contents/{id}/feedback`
- `GET /api/v1/dashboard/recommendations`

## 7. Frontend Information Architecture (MVP)

- `/insights` opportunity board
- `/ideas` idea workspace
- `/studio` trend content studio (platform switch for Xiaohongshu / WeChat)

## 8. Error Handling

- Ingestion failures: mark source status degraded, keep partial pipeline running
- LLM generation failures: return fallback prompt+template suggestions
- Scoring failures: provide draft without scores and explicit warning state
- Feedback anomalies: keep raw values and apply validation gate for tuning

## 9. Testing Strategy (MVP)

- Backend:
  - Unit tests for scoring logic and ranking determinism
  - Contract tests for insight/idea/content/feedback APIs
  - Ingestion pipeline tests for de-dup and cluster stability

- Frontend:
  - Core flow tests for insight -> idea -> studio
  - Rendering/interaction tests for score breakdown and optimization hints

- End-to-end:
  - Seed signals -> generate opportunities -> generate content -> submit feedback

## 10. Risks And Mitigations

- Data source instability: use adapter abstraction and graceful degradation
- Content homogenization: enforce angle diversity constraints in idea generation
- Score drift from real outcomes: introduce periodic backtesting on feedback records
- Multi-objective conflict: expose objective weights in settings (defaults preserved)

## 11. Iteration Roadmap

Phase 1 (MVP):
- Signal ingestion, opportunity board, idea generation, trend draft generation, manual feedback

Phase 2:
- Auto platform metrics sync, smarter ranking calibration, reusable writing style profiles

Phase 3:
- Agentic research and auto A/B suggestion system

## 12. Out Of Scope For MVP

- Full enterprise RBAC and approval chains
- Advanced model-based ranking/learning-to-rank
- Fully automated publishing across all channels

## 13. Approval Record

Approved by user during interactive design review on 2026-03-05.
