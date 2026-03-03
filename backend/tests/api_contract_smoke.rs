use magic_backend::domain::{
    ai_generation::schema::{AddGeneratorReq, SearchGeneratorsReq},
    configuration::schema::{AddSystemConfigReq, SearchProviderConfigsReq},
    content::schema::{AddContentReq, SearchContentReq},
    financing::schema::{CreateBusinessPlanReq, GetBusinessPlanReq},
    identity::schema::{AddUserReq, SearchRolesReq},
    news::schema::{AddNewsReq, SearchNewsReq},
    publishing::schema::{AddPublisherReq, SearchPublishTasksReq},
    scheduling::task::{AddTaskReq, SearchTasksReq},
};

#[test]
fn financing_schema_serde_smoke() {
    let req = CreateBusinessPlanReq {
        title: "t".into(),
        content: "c".into(),
        industry: "i".into(),
        region: "r".into(),
        financing_amount: 1.0,
        company_size: "s".into(),
    };
    let v = serde_json::to_value(req).expect("serialize financing create");
    assert_eq!(v["title"], "t");

    let get_req = GetBusinessPlanReq { id: 42 };
    let v = serde_json::to_value(get_req).expect("serialize financing get");
    assert_eq!(v["id"], 42);
}

#[test]
fn news_schema_serde_smoke() {
    let req = AddNewsReq {
        title: "n".into(),
        content: "body".into(),
        category: Some("cat".into()),
        region: Some("cn".into()),
        industry: None,
        source: Some("src".into()),
    };
    let v = serde_json::to_value(req).expect("serialize news add");
    assert_eq!(v["title"], "n");

    let search = SearchNewsReq {
        query: "q".into(),
        category: None,
        region: None,
        industry: None,
        page: 1,
        limit: 10,
    };
    let v = serde_json::to_value(search).expect("serialize news search");
    assert_eq!(v["query"], "q");
}

#[test]
fn content_schema_serde_smoke() {
    let req = AddContentReq {
        title: "t".into(),
        content: "body".into(),
        content_type: Some("article".into()),
        tags: Some(vec!["rust".into(), "axum".into()]),
    };
    let v = serde_json::to_value(req).expect("serialize content add");
    assert_eq!(v["content_type"], "article");

    let search = SearchContentReq {
        query: Some("rust".into()),
        content_type: None,
        status: None,
        tag: None,
        page: 1,
        limit: 20,
    };
    let v = serde_json::to_value(search).expect("serialize content search");
    assert_eq!(v["limit"], 20);
}

#[test]
fn configuration_schema_serde_smoke() {
    let req = AddSystemConfigReq {
        key: "k".into(),
        value: "v".into(),
        description: Some("d".into()),
        category: Some("cat".into()),
    };
    let v = serde_json::to_value(req).expect("serialize config add");
    assert_eq!(v["key"], "k");

    let search = SearchProviderConfigsReq {
        provider_name: Some("openai".into()),
        config_key: Some("api_key".into()),
        page: 1,
        limit: 10,
    };
    let v = serde_json::to_value(search).expect("serialize provider search");
    assert_eq!(v["provider_name"], "openai");
}

#[test]
fn ai_schema_serde_smoke() {
    let req = AddGeneratorReq {
        name: "g".into(),
        provider: "openai".into(),
        model: "gpt".into(),
        api_key: Some("k".into()),
        api_endpoint: None,
        enabled: Some(true),
    };
    let v = serde_json::to_value(req).expect("serialize ai add generator");
    assert_eq!(v["provider"], "openai");

    let search = SearchGeneratorsReq {
        provider: Some("openai".into()),
        model: None,
        enabled: Some(true),
        page: 1,
        limit: 50,
    };
    let v = serde_json::to_value(search).expect("serialize ai search");
    assert_eq!(v["enabled"], true);
}

#[test]
fn publishing_schema_serde_smoke() {
    let req = AddPublisherReq {
        name: "pub".into(),
        platform: "wechat".into(),
        platform_id: Some("id".into()),
        enabled: Some(true),
    };
    let v = serde_json::to_value(req).expect("serialize publishing add");
    assert_eq!(v["platform"], "wechat");

    let search = SearchPublishTasksReq {
        status: Some("pending".into()),
        platform: None,
        created_at: None,
        page: 1,
        limit: 10,
    };
    let v = serde_json::to_value(search).expect("serialize publishing task search");
    assert_eq!(v["status"], "pending");
}

#[test]
fn scheduling_schema_serde_smoke() {
    let req = AddTaskReq {
        name: "task".into(),
        description: Some("d".into()),
        task_type: "cron".into(),
        scheduler_id: 1,
        cron_expression: Some("* * * * *".into()),
        enabled: Some(true),
    };
    let v = serde_json::to_value(req).expect("serialize scheduling add task");
    assert_eq!(v["scheduler_id"], 1);

    let search = SearchTasksReq {
        name: None,
        task_type: Some("cron".into()),
        scheduler_id: None,
        enabled: Some(true),
        page: 1,
        limit: 10,
    };
    let v = serde_json::to_value(search).expect("serialize scheduling search");
    assert_eq!(v["task_type"], "cron");
}

#[test]
fn identity_schema_serde_smoke() {
    let req = AddUserReq {
        username: "u".into(),
        email: "u@example.com".into(),
        password: Some("p".into()),
        role: Some("admin".into()),
    };
    let v = serde_json::to_value(req).expect("serialize identity add user");
    assert_eq!(v["username"], "u");

    let search = SearchRolesReq {
        name: Some("admin".into()),
        description: None,
        page: 1,
        limit: 10,
    };
    let v = serde_json::to_value(search).expect("serialize identity search roles");
    assert_eq!(v["name"], "admin");
}
