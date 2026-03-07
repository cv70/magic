/**
 * 内容生成领域分类配置
 * 定义了 9 个主领域及其子分类、示例提示词和相关元数据
 */

export interface Domain {
  id: string
  name: string
  icon: string
  description: string
  color: string
  subcategories: string[]
  examplePrompts: string[]
}

export const DOMAINS: Record<string, Domain> = {
  'film-drama': {
    id: 'film-drama',
    name: '影视与短剧制作',
    icon: '🎬',
    description: '生成专业的剧本、分镜、台词和场景描述',
    color: 'red',
    subcategories: ['剧本', '分镜', '台词', '场景描述', '人物小传'],
    examplePrompts: [
      '创作一个20分钟短剧的完整剧本，主题是...',
      '为以下场景写3分钟的精准对话...',
      '设计一部悬疑短剧的分镜台本...'
    ]
  },
  'ai-comics': {
    id: 'ai-comics',
    name: 'AI漫剧/动态漫',
    icon: '🎭',
    description: '生成漫画故事板、对话气泡和人物设定',
    color: 'purple',
    subcategories: ['故事板', '对话设定', '人物设定', '场景背景', '特效描述'],
    examplePrompts: [
      '创作一个6格漫画故事，主题是爱情喜剧...',
      '为以下角色创建详细的人物设定卡...',
      '设计一部动态漫的分镜脚本和对话...'
    ]
  },
  'micro-drama': {
    id: 'micro-drama',
    name: '微短剧（真人/AI混合）',
    icon: '📱',
    description: '针对竖屏的爆点设计、BGM建议和脚本',
    color: 'orange',
    subcategories: ['竖屏脚本', '爆点设计', 'BGM建议', '字幕设计', '转场方案'],
    examplePrompts: [
      '为15秒竖屏短视频设计爆点和转场...',
      '创作一个微短剧的完整脚本和BGM清单...',
      '设计一个搞笑类微短剧的剧情框架...'
    ]
  },
  'animation': {
    id: 'animation',
    name: '动画电影/长视频预演',
    icon: '🎞️',
    description: '角色设定、世界观设定和动作描述',
    color: 'blue',
    subcategories: ['角色设定', '世界观', '动作描述', '视觉风格', '音效设计'],
    examplePrompts: [
      '为一部动画电影创建完整的世界观和背景设定...',
      '设计5个主要角色的详细设定档案...',
      '编写一段复杂的打斗场景的动作描述...'
    ]
  },
  'marketing': {
    id: 'marketing',
    name: '营销与电商广告',
    icon: '📣',
    description: '产品文案、卖点提炼和CTA设计',
    color: 'green',
    subcategories: ['产品文案', '卖点提炼', 'CTA设计', '用户评测文案', '促销文案'],
    examplePrompts: [
      '为这款产品写一段极具吸引力的广告文案...',
      '提炼这个产品的核心卖点和痛点解决方案...',
      '设计一个电商推广页面的完整文案框架...'
    ]
  },
  'education': {
    id: 'education',
    name: '教育与知识出版',
    icon: '📚',
    description: '知识点梳理、测验题目和个性化讲义',
    color: 'cyan',
    subcategories: ['知识点梳理', '测验题目', '讲义内容', '学习路线', '案例分析'],
    examplePrompts: [
      '为一堂课程创建完整的知识体系和讲义...',
      '根据学习目标创建10道测验题及答案解析...',
      '设计一个个性化的学习路线和复习计划...'
    ]
  },
  'interactive': {
    id: 'interactive',
    name: '交互式教材与绘本',
    icon: '🖼️',
    description: '分支剧情、答题逻辑和绘本台词',
    color: 'magenta',
    subcategories: ['分支剧情', '答题逻辑', '绘本台词', '交互设计', '游戏化元素'],
    examplePrompts: [
      '创建一个互动式故事，包含3个关键决策点...',
      '为儿童绘本设计故事内容和配套问答...',
      '设计一个教育性游戏的完整交互流程...'
    ]
  },
  'news-media': {
    id: 'news-media',
    name: '新闻媒体与广电',
    icon: '📰',
    description: '新闻稿、播报稿和融媒体素材',
    color: 'volcano',
    subcategories: ['新闻稿', '播报稿', '融媒体素材', '字幕稿', '评论文章'],
    examplePrompts: [
      '根据事件摘要写一篇500字的新闻稿...',
      '为电视新闻节目编写专业的播报稿...',
      '创建一条融媒体新闻的完整内容套装...'
    ]
  },
  'music-audio': {
    id: 'music-audio',
    name: '音乐与音频娱乐',
    icon: '🎵',
    description: '歌词创作、音频脚本和播客提纲',
    color: 'gold',
    subcategories: ['歌词创作', '音频脚本', '播客提纲', '配音文案', '音效描述'],
    examplePrompts: [
      '创作一首主题为"青春梦想"的流行歌曲歌词...',
      '编写一期播客节目的完整提纲和脚本...',
      '为有声书章节编写专业的配音脚本...'
    ]
  }
}

/**
 * 快速获取所有领域列表（按 ID）
 */
export const getDomainsList = (): Domain[] => {
  return Object.values(DOMAINS)
}

/**
 * 快速获取领域 ID 列表
 */
export const getDomainIds = (): string[] => {
  return Object.keys(DOMAINS)
}

/**
 * 根据 ID 获取领域信息
 */
export const getDomainById = (id: string): Domain | undefined => {
  return DOMAINS[id]
}

/**
 * 根据领域 ID 获取领域名称
 */
export const getDomainName = (id: string): string => {
  return DOMAINS[id]?.name || '未知领域'
}

/**
 * 根据领域 ID 获取图标
 */
export const getDomainIcon = (id: string): string => {
  return DOMAINS[id]?.icon || '❓'
}
