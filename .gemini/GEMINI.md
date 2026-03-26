# Gemini Configuration for Workspace Knowledge Management

This file outlines how Gemini utilizes its Knowledge Graph (Memory MCP) for enhanced understanding and interaction within the workspace.

## Knowledge Graph (Memory MCP) Usage

The Memory Multi-Capability Provider (MCP) enables Gemini to store, retrieve, and reason about structured information in the form of a knowledge graph. This capability is crucial for maintaining context, understanding project nuances, and providing more intelligent assistance in future development tasks.

### Core Concepts

The knowledge graph is built upon three fundamental components:

1.  **Entities:** Represent distinct objects, concepts, or components within the project.
    *   **Examples:** `User Model`, `Auth Service`, `API Endpoint /users`, `Database Table 'products'`, `React Component 'UserProfile'`.
    *   **Attributes:** Each entity has a `name` (unique identifier) and an `entityType` (category).

2.  **Relations:** Describe how entities are connected to each other. Relations are directional and represent a specific type of interaction or dependency.
    *   **Examples:** `Auth Service` --(calls)--> `User Model`, `API Endpoint /users` --(uses)--> `Auth Service`, `React Component 'UserProfile'` --(renders)--> `User Model data`.
    *   **Attributes:** Each relation has a `from` entity, a `to` entity, and a `relationType`.

3.  **Observations:** Provide additional, descriptive information or attributes about an entity. These are free-form text snippets that add context.
    *   **Examples:** For `User Model`: "Contains fields: `username`, `email`, `passwordHash`", "Managed by the Backend Team", "Uses Mongoose schema validation".
    *   **Attributes:** Each observation is linked to an `entityName` and contains `contents` (an array of strings).

### Usage Guidelines for Gemini

To effectively leverage the knowledge graph, Gemini will adhere to the following guidelines:

*   **Proactive Knowledge Capture:** Gemini will proactively identify and store key entities, relationships, and observations from code, documentation, and user interactions that are relevant for future tasks.
*   **Contextual Querying:** Before making significant changes or providing complex explanations, Gemini will query the knowledge graph to gather relevant context, understand dependencies, and identify potential impacts.
*   **Maintaining Accuracy:** Gemini will strive to keep the knowledge graph up-to-date. If code changes invalidate existing knowledge, Gemini will propose updates to the graph.
*   **Prioritizing User-Provided Information:** Explicit instructions or facts provided by the user will take precedence and be stored in the graph to personalize interactions.

### Maintenance and Evolution

The knowledge graph is a living document that evolves with the project. Users can also directly interact with the graph to refine its content:

*   **Adding Knowledge:**                                                                                                                              
    ```python                                                                                                                                          
    # Example: Adding a new entity and an observation                                                                                                  
    print(default_api.create_entities(entities=[                                                                                                       
        {"name": "Payment Gateway", "entityType": "External Service", "observations": ["Integrates with Razorpay", "Handles all transaction processing"]}                                                                                                                                          
    ]))                                                                                                                                                
    ```                                                                                                                                                
    ```python                                                                                                                                          
    # Example: Adding a relation                                                                                                                       
    print(default_api.create_relations(relations=[                                                                                                     
        {"from": "Order Service", "to": "Payment Gateway", "relationType": "initiates transaction with"}                                               
    ]))                                                                                                                                                
    ```                                                                                                                                                
*   **Querying Knowledge:**                                                                                                                            
    ```python                                                                                                                                          
    # Example: Reading the entire graph                                                                                                                
    print(default_api.read_graph())                                                                                                                    
    ```                                                                                                                                                
    ```python                                                                                                                                          
    # Example: Searching for entities related to "authentication"                                                                                      
    print(default_api.search_nodes(query="authentication"))                                                                                            
    ```                                                                                                                                                
*   **Updating/Deleting Knowledge:**                                                                                                                   
    *   If an entity or relation becomes obsolete, it should be removed to maintain graph integrity.                                                   
    *   Gemini will propose deletions or updates when it detects inconsistencies.                                                                      
    ```python                                                                                                                                          
    # Example: Deleting an entity                                                                                                                      
    print(default_api.delete_entities(entityNames=["Old Feature Module"]))                                                                             
    ``` 

By actively maintaining and utilizing this knowledge graph, Gemini can provide more informed, accurate, and efficient assistance throughout the development lifecycle.

## MCP Server Configurations

This section details the Multi-Capability Provider (MCP) servers configured for Gemini, along with the tools they provide.

### 🟢 memory - Ready (9 tools)
  Tools:
  - add_observations
    Add new observations to existing entities in the knowledge graph
  - create_entities
    Create multiple new entities in the knowledge graph
  - create_relations
    Create multiple new relations between entities in the knowledge graph. Relations should be in active voice
  - delete_entities
    Delete multiple entities and their associated relations from the knowledge graph
  - delete_observations
    Delete specific observations from entities in the knowledge graph
  - delete_relations
    Delete multiple relations from the knowledge graph
  - open_nodes
    Open specific nodes in the knowledge graph by their names
  - read_graph
    Read the entire knowledge graph
  - search_nodes
    Search for nodes in the knowledge graph based on a query

### 🟢 sequential-thinking - Ready (1 tool)
  Tools:
  - sequentialthinking
    A detailed tool for dynamic and reflective problem-solving through thoughts.
    This tool helps analyze problems through a flexible thinking process that can adapt and evolve.
    Each thought can build on, question, or revise previous insights as understanding deepens.

    When to use this tool:
    - Breaking down complex problems into steps
    - Planning and design with room for revision
    - Analysis that might need course correction
    - Problems where the full scope might not be clear initially
    - Problems that require a multi-step solution
    - Tasks that need to maintain context over multiple steps
    - Situations where irrelevant information needs to be filtered out

    Key features:
    - You can adjust total_thoughts up or down as you progress
    - You can question or revise previous thoughts
    - You can add more thoughts even after reaching what seemed like the end
    - You can express uncertainty and explore alternative approaches
    - Not every thought needs to build linearly - you can branch or backtrack
    - Generates a solution hypothesis
    - Verifies the hypothesis based on the Chain of Thought steps
    - Repeats the process until satisfied
    - Provides a correct answer

    Parameters explained:
    - thought: Your current thinking step, which can include:
    * Regular analytical steps
    * Revisions of previous thoughts
    * Questions about previous decisions
    * Realizations about needing more analysis
    * Changes in approach
    * Hypothesis generation
    * Hypothesis verification
    - next_thought_needed: True if you need more thinking, even if at what seemed like the end
    - thought_number: Current number in sequence (can go beyond initial total if needed)
    - total_thoughts: Current estimate of thoughts needed (can be adjusted up/down)
    - is_revision: A boolean indicating if this thought revises previous thinking
    - revises_thought: If is_revision is true, which thought number is being reconsidered
    - branch_from_thought: If branching, which thought number is the branching point
    - branch_id: Identifier for the current branch (if any)
    - needs_more_thoughts: If reaching end but realizing more thoughts needed

    You should:
    1. Start with an initial estimate of needed thoughts, but be ready to adjust
    2. Feel free to question or revise previous thoughts
    3. Don't hesitate to add more thoughts if needed, even at the "end"
    4. Express uncertainty when present
    5. Mark thoughts that revise previous thinking or branch into new paths
    6. Ignore information that is irrelevant to the current step
    7. Generate a solution hypothesis when appropriate
    8. Verify the hypothesis based on the Chain of Thought steps
    9. Repeat the process until satisfied with the solution
    10. Provide a single, ideally correct answer as the final output
    11. Only set next_thought_needed to false when truly done and a satisfactory answer is reached