/**
 * 每个领域对应的专属模板和参数配置
 * 用于 AI 生成时自动加载领域相关的提示词模板和参数字段
 */

export interface DomainTemplate {
  id: string
  domainId: string
  name: string
  description: string
  basePrompt: string
  paramFields: ParamField[]
  exampleOutput: string
}

export interface ParamField {
  key: string
  label: string
  type: 'text' | 'textarea' | 'select' | 'number' | 'slider'
  placeholder?: string
  required?: boolean
  options?: Array<{ label: string; value: string | number }>
  min?: number
  max?: number
  defaultValue?: string | number
}

export const DOMAIN_TEMPLATES: Record<string, DomainTemplate[]> = {
  'film-drama': [
    {
      id: 'film-drama-script',
      domainId: 'film-drama',
      name: '短剧剧本生成',
      description: '生成完整的短剧剧本，包含场景、角色和对话',
      basePrompt: '你是一位专业编剧。请根据以下参数创作一个短剧剧本：\n主题：{theme}\n时长：{duration}分钟\n场景：{setting}\n风格：{style}\n目标观众：{audience}',
      paramFields: [
        {
          key: 'theme',
          label: '故事主题',
          type: 'text',
          placeholder: '例如：爱情、悬疑、喜剧...',
          required: true
        },
        {
          key: 'duration',
          label: '剧本时长',
          type: 'slider',
          min: 5,
          max: 60,
          defaultValue: 20
        },
        {
          key: 'setting',
          label: '主要场景',
          type: 'text',
          placeholder: '例如：现代都市、古代宫廷...',
          required: true
        },
        {
          key: 'style',
          label: '创作风格',
          type: 'select',
          options: [
            { label: '温情励志', value: 'heartwarming' },
            { label: '悬疑惊悚', value: 'suspense' },
            { label: '爆笑喜剧', value: 'comedy' },
            { label: '热血青春', value: 'youth' }
          ],
          required: true
        },
        {
          key: 'audience',
          label: '目标观众',
          type: 'select',
          options: [
            { label: '全年龄', value: 'all-ages' },
            { label: '青少年', value: 'teens' },
            { label: '成人', value: 'adults' }
          ],
          required: true
        }
      ],
      exampleOutput: '第一场 [办公室 - 上午]\n角色：李明、王芳\n[李明坐在办公桌前，面容疲惫]\n李明：又是加班到深夜...\n[王芳敲门进入]\n王芳：你又在这儿？...'
    },
    {
      id: 'film-drama-dialogue',
      domainId: 'film-drama',
      name: '精准对话生成',
      description: '为特定场景生成角色间的自然流畅的对话',
      basePrompt: '请为以下场景创作{count}句自然的对话：\n场景：{scene}\n角色关系：{relationship}\n情感基调：{emotion}',
      paramFields: [
        {
          key: 'scene',
          label: '场景描述',
          type: 'textarea',
          placeholder: '描述场景的时间、地点、背景...',
          required: true
        },
        {
          key: 'relationship',
          label: '角色关系',
          type: 'text',
          placeholder: '例如：男朋友和女朋友、医生和患者...',
          required: true
        },
        {
          key: 'emotion',
          label: '情感基调',
          type: 'select',
          options: [
            { label: '温情', value: 'tender' },
            { label: '紧张', value: 'tense' },
            { label: '轻松', value: 'relaxed' },
            { label: '悲伤', value: 'sad' }
          ],
          required: true
        },
        {
          key: 'count',
          label: '对话轮数',
          type: 'slider',
          min: 3,
          max: 20,
          defaultValue: 10
        }
      ],
      exampleOutput: '人物A：你为什么一直在看手机？\n人物B：抱歉，我只是...有点担心。\n人物A：担心什么？告诉我。'
    }
  ],

  'ai-comics': [
    {
      id: 'ai-comics-storyboard',
      domainId: 'ai-comics',
      name: '漫画故事板生成',
      description: '生成连贯的漫画故事板，包含每格的场景、对话和动作',
      basePrompt: '你是一位资深漫画编剧。请创作一个{panels}格的漫画故事：\n主题：{theme}\n风格：{style}\n每格需要包含：场景描述、角色对话、动作指示',
      paramFields: [
        {
          key: 'theme',
          label: '故事主题',
          type: 'text',
          placeholder: '例如：日常搞笑、冒险奇遇...',
          required: true
        },
        {
          key: 'panels',
          label: '故事格数',
          type: 'slider',
          min: 4,
          max: 16,
          defaultValue: 6
        },
        {
          key: 'style',
          label: '漫画风格',
          type: 'select',
          options: [
            { label: '日式', value: 'japanese' },
            { label: '美式', value: 'american' },
            { label: '港漫', value: 'hongkong' },
            { label: '欧洲', value: 'european' }
          ],
          required: true
        }
      ],
      exampleOutput: '第1格：[明亮的办公室背景] 角色A坐在办公桌前，表情认真。\n对话：A说"又是报告截止日期了！"\n动作：A敲打键盘，汗滴掉落\n\n第2格：...'
    }
  ],

  'micro-drama': [
    {
      id: 'micro-drama-script',
      domainId: 'micro-drama',
      name: '竖屏微短剧脚本',
      description: '为竖屏视频生成爆点突出、节奏快的脚本',
      basePrompt: '请创作一个{duration}秒的竖屏短视频脚本：\n类型：{type}\n爆点主题：{climax_theme}\n目标情绪：{emotion}\n要求：节奏快、转场多、有强烈的爆点',
      paramFields: [
        {
          key: 'duration',
          label: '视频时长',
          type: 'select',
          options: [
            { label: '15秒', value: 15 },
            { label: '30秒', value: 30 },
            { label: '60秒', value: 60 }
          ],
          required: true
        },
        {
          key: 'type',
          label: '视频类型',
          type: 'select',
          options: [
            { label: '搞笑', value: 'comedy' },
            { label: '温情', value: 'emotional' },
            { label: '惊悚', value: 'thriller' },
            { label: '励志', value: 'inspirational' }
          ],
          required: true
        },
        {
          key: 'climax_theme',
          label: '爆点主题',
          type: 'text',
          placeholder: '例如：反转、亲情、笑点...',
          required: true
        },
        {
          key: 'emotion',
          label: '目标情绪',
          type: 'text',
          placeholder: '例如：惊喜、感动、大笑...',
          required: true
        }
      ],
      exampleOutput: '[0-5秒] 场景：餐厅。女孩拿起菜单，表情认真。\n对话：F："请给我来一份..."\n[5-10秒] 转场：快速翻转镜头\n...'
    }
  ],

  'animation': [
    {
      id: 'animation-character',
      domainId: 'animation',
      name: '角色设定生成',
      description: '创建动画电影的主要角色设定档案',
      basePrompt: '请为一部动画电影创建{count}个主要角色的详细设定：\n故事背景：{world_setting}\n时代背景：{time_period}\n需要包括：外形、性格、背景故事、能力/特点',
      paramFields: [
        {
          key: 'count',
          label: '角色数量',
          type: 'slider',
          min: 1,
          max: 10,
          defaultValue: 5
        },
        {
          key: 'world_setting',
          label: '世界观设定',
          type: 'textarea',
          placeholder: '描述故事发生的世界...',
          required: true
        },
        {
          key: 'time_period',
          label: '时代背景',
          type: 'select',
          options: [
            { label: '现代', value: 'modern' },
            { label: '古代', value: 'ancient' },
            { label: '未来', value: 'future' },
            { label: '奇幻', value: 'fantasy' }
          ],
          required: true
        }
      ],
      exampleOutput: '【角色1】\n名字：林晓雨\n年龄：18岁\n职业：魔法学徒\n外形：长黑发，眼睛有特殊的魔力闪光\n性格：温柔但坚定，充满冒险精神\n...'
    }
  ],

  'marketing': [
    {
      id: 'marketing-product-copy',
      domainId: 'marketing',
      name: '产品文案生成',
      description: '为产品生成吸引力强的营销文案',
      basePrompt: '请为以下产品创作一段{length}字的营销文案：\n产品名称：{product_name}\n核心卖点：{core_benefits}\n目标人群：{target_audience}\n风格：{style}',
      paramFields: [
        {
          key: 'product_name',
          label: '产品名称',
          type: 'text',
          placeholder: '例如：智能空气净化器...',
          required: true
        },
        {
          key: 'core_benefits',
          label: '核心卖点',
          type: 'textarea',
          placeholder: '列举产品的主要优势和特点',
          required: true
        },
        {
          key: 'target_audience',
          label: '目标人群',
          type: 'text',
          placeholder: '例如：年轻上班族、新手父母...',
          required: true
        },
        {
          key: 'style',
          label: '文案风格',
          type: 'select',
          options: [
            { label: '高端大气', value: 'premium' },
            { label: '幽默亲切', value: 'humorous' },
            { label: '紧迫煽情', value: 'urgent' },
            { label: '专业技术', value: 'technical' }
          ],
          required: true
        },
        {
          key: 'length',
          label: '文案长度',
          type: 'slider',
          min: 50,
          max: 500,
          defaultValue: 200
        }
      ],
      exampleOutput: '【产品文案】\n这不仅仅是一个空气净化器，它是你和家人的健康守护者。采用日本进口 HEPA 滤芯，...'
    }
  ],

  'education': [
    {
      id: 'education-courseware',
      domainId: 'education',
      name: '课程讲义生成',
      description: '为课程创建完整的知识体系和讲义',
      basePrompt: '请为一堂课程创建完整的讲义和知识体系：\n课程主题：{topic}\n学生水平：{level}\n课时数：{hours}\n需要包括：大纲、知识点、案例、小结',
      paramFields: [
        {
          key: 'topic',
          label: '课程主题',
          type: 'text',
          placeholder: '例如：Python 基础、经济学原理...',
          required: true
        },
        {
          key: 'level',
          label: '学生水平',
          type: 'select',
          options: [
            { label: '初级入门', value: 'beginner' },
            { label: '中级进阶', value: 'intermediate' },
            { label: '高级专家', value: 'advanced' }
          ],
          required: true
        },
        {
          key: 'hours',
          label: '课时数',
          type: 'slider',
          min: 1,
          max: 40,
          defaultValue: 5
        }
      ],
      exampleOutput: '【第一章】基础概念\n\n学习目标：\n- 理解...的定义\n- 掌握...的应用场景\n\n知识点 1.1：定义\n...'
    }
  ],

  'interactive': [
    {
      id: 'interactive-story',
      domainId: 'interactive',
      name: '互动故事生成',
      description: '创建包含多分支决策的互动式故事',
      basePrompt: '请创作一个互动式故事，包含{branches}个关键决策点：\n主题：{theme}\n目标读者：{audience}\n每个决策点需要有至少 2 个选项，并导向不同的故事发展',
      paramFields: [
        {
          key: 'theme',
          label: '故事主题',
          type: 'text',
          placeholder: '例如：冒险、悬疑、爱情...',
          required: true
        },
        {
          key: 'branches',
          label: '决策点数',
          type: 'slider',
          min: 3,
          max: 10,
          defaultValue: 5
        },
        {
          key: 'audience',
          label: '目标读者',
          type: 'select',
          options: [
            { label: '儿童', value: 'children' },
            { label: '青少年', value: 'teens' },
            { label: '成人', value: 'adults' }
          ],
          required: true
        }
      ],
      exampleOutput: '你走进了一间古老的书店...\n\n【决策点 1】\n你应该：\nA. 走向左边的神秘书架\nB. 去找店主询问\n\n如果选择 A：故事继续...'
    }
  ],

  'news-media': [
    {
      id: 'news-media-article',
      domainId: 'news-media',
      name: '新闻稿生成',
      description: '根据事件信息生成专业的新闻稿',
      basePrompt: '请根据以下信息撰写一篇{style}新闻稿（约{length}字）：\n事件名称：{event_name}\n关键信息：{key_info}\n发生地点：{location}\n发生时间：{time}',
      paramFields: [
        {
          key: 'event_name',
          label: '事件标题',
          type: 'text',
          placeholder: '例如：某企业推出新产品...',
          required: true
        },
        {
          key: 'key_info',
          label: '核心信息',
          type: 'textarea',
          placeholder: '描述事件的重要细节',
          required: true
        },
        {
          key: 'location',
          label: '发生地点',
          type: 'text',
          placeholder: '例如：北京国家会议中心...',
          required: true
        },
        {
          key: 'time',
          label: '发生时间',
          type: 'text',
          placeholder: '例如：2024 年 3 月 7 日...',
          required: true
        },
        {
          key: 'style',
          label: '新闻风格',
          type: 'select',
          options: [
            { label: '新闻通稿', value: 'press_release' },
            { label: '深度报道', value: 'feature' },
            { label: '快讯', value: 'breaking' }
          ],
          required: true
        },
        {
          key: 'length',
          label: '文章长度',
          type: 'slider',
          min: 300,
          max: 1500,
          defaultValue: 500
        }
      ],
      exampleOutput: '【新闻标题】\n\n【导语】...\n\n【新闻正文】...\n\n【相关链接】...'
    }
  ],

  'music-audio': [
    {
      id: 'music-audio-lyrics',
      domainId: 'music-audio',
      name: '歌词创作',
      description: '为歌曲创作富有表现力的歌词',
      basePrompt: '请创作一首歌词，要求：\n歌曲名称：{song_title}\n主题：{theme}\n风格：{style}\n结构：{structure}\n目标受众：{audience}',
      paramFields: [
        {
          key: 'song_title',
          label: '歌曲名称',
          type: 'text',
          placeholder: '例如：青春梦想...',
          required: true
        },
        {
          key: 'theme',
          label: '歌曲主题',
          type: 'textarea',
          placeholder: '描述歌曲要表达的情感和意境',
          required: true
        },
        {
          key: 'style',
          label: '音乐风格',
          type: 'select',
          options: [
            { label: '流行', value: 'pop' },
            { label: '摇滚', value: 'rock' },
            { label: '民谣', value: 'folk' },
            { label: '说唱', value: 'hiphop' }
          ],
          required: true
        },
        {
          key: 'structure',
          label: '歌曲结构',
          type: 'select',
          options: [
            { label: '副歌*2+主歌*2', value: 'standard' },
            { label: '主歌*3+副歌*2', value: 'verse_heavy' },
            { label: '自由形式', value: 'free' }
          ],
          required: true
        },
        {
          key: 'audience',
          label: '目标受众',
          type: 'text',
          placeholder: '例如：90 后青年...',
          required: true
        }
      ],
      exampleOutput: '【副歌】\n在梦想的路上，我从不曾害怕\n即使摔跤也要爬起来\n...\n\n【主歌】\n月亮见证我的努力...'
    }
  ]
}

/**
 * 根据领域 ID 获取该领域的所有模板
 */
export const getTemplatesByDomain = (domainId: string): DomainTemplate[] => {
  return DOMAIN_TEMPLATES[domainId] || []
}

/**
 * 根据模板 ID 获取具体模板
 */
export const getTemplateById = (templateId: string): DomainTemplate | undefined => {
  for (const templates of Object.values(DOMAIN_TEMPLATES)) {
    const found = templates.find(t => t.id === templateId)
    if (found) return found
  }
  return undefined
}

/**
 * 根据参数值构建最终的生成提示词
 */
export const buildPromptFromTemplate = (
  template: DomainTemplate,
  params: Record<string, any>
): string => {
  let prompt = template.basePrompt
  for (const [key, value] of Object.entries(params)) {
    prompt = prompt.replace(`{${key}}`, String(value))
  }
  return prompt
}
